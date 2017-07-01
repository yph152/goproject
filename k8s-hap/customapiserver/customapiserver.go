package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/emicklei/go-restful"
	"k8s.io/k8s-hap/client/etcdclient"
	"k8s.io/k8s-hap/k8sv1beta1"
	"net/http"
)

/*type K8sScaleSpec struct {
	Replicas          int32  `json:"replicas,omitempty"`
	MinReplicas       int32  `json:"minReplicas,omitempty"`
	MaxReplicas       int32  `json:"maxReplicas,omitempty"`
	Name              string `json:"name,omitempty"`
	Namespace         string `json:"namespace,omitemptyi"`
	TargetPercentage  int32  `json:"targetPercentage,omitempty"`
	CurrentPercentage int32  `json:"currentPercentage,omitempty"`
	CurrentReplicas   int32  `json:"currentReplicas,omitempty"`
}*/

type K8sScaleResource struct {
	K8sSpec    *k8sv1beta1.K8sScaleSpec
	EtcdClient *etcdclient.EtcdClient
}

func init() {

}

func (u K8sScaleResource) Register(container *restful.Container) {
	//创建新的WebService
	ws := new(restful.WebService)

	//设定WebService对应的路径("/users")和支持的MIME类型(restful.MIME_XML/restful.MIME_JSON)
	ws.
		Path("/api/v1").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	//添加路由: POST /{namespace} --> u.createK8sScalerSpec
	ws.Route(ws.POST("/autoscalers/{namespace}/{name}").To(u.createK8sScalerSpec))
	//添加路由: DELETE /{namespace}/{name} --> u.deleteK8sScalerSpec
	ws.Route(ws.DELETE("/autoscalers/{namespace}/{name}").To(u.deleteK8sScalerSpec))
	//添加路由: GET /{namespace}/{name} --> u.getK8sScalerSpec
	ws.Route(ws.GET("/autoscalers/{namespace}/{name}").To(u.getK8sScalerSpec))
	//添加路由: GET /{user-id} --> u.getK8sScalerSpec
	ws.Route(ws.PUT("/autoscalers/{namespace}").To(u.updateK8sScalerSpec))
	//添加路由: GET /{namespace} --> u.showNamespace
	ws.Route(ws.GET("/namespaces").To(u.showNamespace))
	//添加路由: GET/{namespace} --> u.showAutoScalers
	ws.Route(ws.GET("/autoscalers/{namespace}").To(u.showAutoScalers))
	//添加路由: POST/namespaces/{namespace}
	ws.Route(ws.POST("/namespaces/{namespace}").To(u.createK8sNamespace))

	//将初始化好的WebService添加到Container中
	container.Add(ws)
}

func (u *K8sScaleResource) createK8sNamespace(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	key := "/api/v1/namespaces/" + namespace

	var usr string
	err := request.ReadEntity(&usr)

	if err != nil {
		u.EtcdClient.EtcdSet(key, namespace)
		response.WriteHeader(http.StatusCreated)
		response.WriteEntity(usr)
	} else {
		response.AddHeader("Content-Type", "application/json")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}
func (u *K8sScaleResource) getK8sScalerSpec(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")

	key := "/api/v1/autoscalers/" + namespace + "/" + name

	value, err := u.EtcdClient.EtcdGet(key)
	if err != nil {
		response.AddHeader("Content-Type", "application/json")
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
		return
	}

	var data k8sv1beta1.K8sScaleSpec

	json.Unmarshal([]byte(value), &data)

	if len(data.Name) == 0 {
		response.AddHeader("Content-Type", "application/json")
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
	} else {
		response.WriteEntity(data)
	}
}

func (u *K8sScaleResource) updateK8sScalerSpec(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")

	key := "/api/v1/autoscalers/" + namespace + "/" + name

	usr := new(k8sv1beta1.K8sScaleSpec)
	err := request.ReadEntity(&usr)
	value, err := json.Marshal(usr)

	if err == nil {
		u.EtcdClient.EtcdUpdate(key, string(value))
		response.WriteEntity(usr)
	} else {
		response.AddHeader("Content-Type", "application/json")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (u *K8sScaleResource) deleteK8sScalerSpec(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")

	key := "/api/v1/autoscalers/" + namespace + "/" + name

	u.EtcdClient.EtcdDelete(key)
}

func (u *K8sScaleResource) createK8sScalerSpec(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")

	key := "/api/v1/autoscalers/" + namespace + "/" + name

	usr := k8sv1beta1.K8sScaleSpec{Namespace: request.PathParameter("namespace"), Name: request.PathParameter("name")}
	err := request.ReadEntity(&usr)

	value, _ := json.Marshal(usr)
	if err == nil {
		u.EtcdClient.EtcdSet(key, string(value))
		response.WriteHeader(http.StatusCreated)
		response.WriteEntity(usr)
	} else {
		response.AddHeader("Content-Type", "application/json")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (u *K8sScaleResource) showNamespace(request *restful.Request, response *restful.Response) {
	//	namespace := request.PathParameter("namespace")

	key := "/api/v1/namespaces"
	/*	for k, _ := range usr {
			test = k
		}
	*/

	namespace, _ := u.EtcdClient.EtcdList(key)
	if len(namespace[0]) == 0 {
		response.AddHeader("Content-Type", "application/json")
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
	} else {
		response.WriteEntity(namespace)
	}

}

func (u *K8sScaleResource) showAutoScalers(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	//	name := request.PathParameter("name")
	key := "/api/v1/autoscalers/" + namespace

	name, _ := u.EtcdClient.EtcdList(key)
	if len(name[0]) == 0 {
		response.AddHeader("Content-Type", "applicatiuon/json")
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
	} else {
		response.WriteEntity(name)
	}
}

var (
	etcdconfig = flag.String("etcdendpoints", "http://127.0.0.1:2379", "endpoints to the etcd")
)

func main() {
	//创建一个空的container

	wsContainer := restful.NewContainer()

	//设定路由为CurlyRoute(快速路由)
	wsContainer.Router(restful.CurlyRouter{})

	flag.Parse()

	str1 := make([]string, 3)
	str1[0] = *etcdconfig
	client, err := etcdclient.NewClient(str1)

	if err != nil {
		fmt.Println(err)
	}

	c := etcdclient.EtcdClient{}
	c.Client = client

	//创建自定义的Resource Handle(此处为K8sScaleResource)
	u := K8sScaleResource{K8sSpec: &k8sv1beta1.K8sScaleSpec{}, EtcdClient: &c}

	//创建WebService,bin并将WebService 加入到Container中

	u.Register(wsContainer)

	fmt.Println("start listening on localhost:8088")
	server := &http.Server{Addr: ":8088", Handler: wsContainer}

	//启动服务
	server.ListenAndServe()
}
