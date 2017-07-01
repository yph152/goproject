package podautoscaler

import (
	"encoding/json"
	//"flag"
	"fmt"
	"k8s.io/client-go/pkg/util/clock"
	"k8s.io/k8s-hap/client/customclient"
	"k8s.io/k8s-hap/client/k8sapiserver"
	"k8s.io/k8s-hap/client/promclient"
	"k8s.io/k8s-hap/k8sv1beta1"
	"k8s.io/k8s-hap/options"
	"k8s.io/k8s-hap/podautoscaler/queue"
	"math"
	"net/http"
	"time"
)

type Autoscaler struct {
	CustomClient *customclient.CustomClient
	K8sClient    *k8sapiserver.K8sapiserver
	PromClient   *promclient.PrometheusClient
}
type Horization struct {
	Auto        *Autoscaler
	QueueScaler *queue.Queue
	QueueWatch  *queue.Queue
}
type HorizationResource struct {
	Horiza  *Horization
	K8sSpec k8sv1beta1.K8sScaleSpec
}

func NewAutoScaler(c *options.AutoscalerConfig) *Autoscaler {
	var autoscaler *Autoscaler
	autoscaler = &Autoscaler{
		CustomClient: &customclient.CustomClient{Client: &http.Client{}, Endpoints: c.CustomClient},
		K8sClient:    k8sapiserver.NewClient(c.Kubeconfig),
		PromClient:   promclient.NewPrometheusClient(c.PrometheusIPPort),
	}
	return autoscaler
}
func HorizationInit(c *options.AutoscalerConfig) *Horization {
	scale := NewAutoScaler(c)
	quescaler := queue.NewQueue(Worker, 10)

	quewatch := queue.NewQueue(Watcher, 20)

	horization := &Horization{Auto: scale, QueueScaler: quescaler, QueueWatch: quewatch}
	return horization
}

const (
	defaultTargetPercentage   = 80
	defaultMinTargetPerentage = 60
)

var downscaleForbiddenWindow = 5 * time.Minute
var upscaleForbiddenWindow = 1 * time.Minute

/*func (h *Horization) CalculateScaleUpLimit(currentReplicas int32, Replicas int32) int32 {
	return int32(math.Max(currentReplicas, Replicas))
}*/

func (h *Horization) Run() {
	kclock := clock.RealClock{}
	ticker := kclock.Tick(15)
	h.Informer()

	for {
		select {
		case <-ticker:
			h.Informer()
		}
	}
}

func (h *Horization) Informer() {
	path := h.Auto.CustomClient.Path("", "")
	url := h.Auto.CustomClient.Url(path)
	namespaces := h.Auto.CustomClient.CustomGetNameSpace(url)
	for _, value := range namespaces {
		path1 := h.Auto.CustomClient.Path(value, "")
		url1 := h.Auto.CustomClient.Url(path1)
		name := h.Auto.CustomClient.CustomGetAutoScalers(url1)
		for _, value1 := range name {
			path2 := h.Auto.CustomClient.Path(value, value1)
			url2 := h.Auto.CustomClient.Url(path2)
			spec := h.Auto.CustomClient.CustomGet(url2)
			h.QueueWatch.Push(spec)

		}
	}
}
func Watcher(spec interface{}) {
	switch v := spec.(type) {
	case HorizationResource:
		WatcherRun(v)
	default:
		fmt.Println("unknow watcher type")
	}
}
func WatcherRun(spec HorizationResource) {
	var tag bool
	var replicas k8sv1beta1.K8sScaleSpec
	if spec.K8sSpec.ForCPUUtilization == true {
		replicas, tag = spec.Horiza.ComputeReplicasForCPUUtilization(spec.K8sSpec)
	} else if spec.K8sSpec.ForMemoryUtilization == true {
		replicas, tag = spec.Horiza.ComputeReplicasForMemoryUtilization(spec.K8sSpec)
	}

	if tag == false {
		if spec.K8sSpec.CurrentReplicas < spec.K8sSpec.MinReplicas {
			spec.K8sSpec.Replicas = replicas.Replicas
			spec.Horiza.QueueScaler.Push(spec)
		}
	} else {
		if spec.K8sSpec.CurrentReplicas < replicas.Replicas {
			spec.K8sSpec.Replicas = replicas.Replicas
			spec.Horiza.QueueScaler.Push(spec)
		}
	}
}
func Worker(spec interface{}) {

	switch v := spec.(type) {
	case HorizationResource:
		WorkerRun(v)
	default:
		fmt.Println("unknow worker type")
	}

}
func WorkerRun(spec HorizationResource) {
	//转换为json数据
	jsonstr, _ := json.Marshal(spec.K8sSpec)
	path := spec.Horiza.Auto.CustomClient.Path(spec.K8sSpec.Namespace, spec.K8sSpec.Name)
	url := spec.Horiza.Auto.CustomClient.Url(path)
	spec.Horiza.Auto.CustomClient.CustomPost(url, jsonstr)
	spec.Horiza.Auto.K8sClient.K8sScaler(spec.K8sSpec.Namespace, spec.K8sSpec.Name, spec.K8sSpec.Replicas)
}
func (h *Horization) ComputeReplicasForCPUUtilization(spec k8sv1beta1.K8sScaleSpec) (k8sv1beta1.K8sScaleSpec, bool) {
	var limit float64 = h.ComputeCPULimit(spec.Namespace, spec.Name)
	var usage float64 = h.ComputeCPUUsage(spec.Namespace, spec.Name)
	//var replicas int32
	var tag bool
	spec.CurrentPercentage = usage / limit
	spec.CurrentReplicas, _ = h.Auto.K8sClient.K8sGetReplicas(spec.Namespace, spec.Name)
	spec.Replicas, _ = h.Auto.K8sClient.K8sGetReplicas(spec.Namespace, spec.Name)
	if spec.CurrentPercentage > spec.TargetPercentage {
		podlimit := spec.CurrentPercentage - spec.TargetPercentage
		podcurrentlimit := h.ComputePodCpuLimit(spec.Namespace, spec.Name)
		replicas := podlimit / podcurrentlimit
		spec.Replicas = int32(math.Ceil(replicas))
		tag = true
		return spec, tag
	} else {
		tag = false
		return spec, tag
	}

}
func (h *Horization) ComputeCPULimit(namespace string, name string) float64 {
	queryExp := string(`sum(rate(kube_pod_container_resource_requests_cpu_cores{namespace="`) + namespace + string(`",container="`) + name + string(`"}[2m]))`)
	metrics := h.Auto.PromClient.GetPrometheus(string(queryExp))
	return metrics[0].Float()

}
func (h *Horization) ComputeCPUUsage(namespace string, name string) float64 {
	queryExp := string(`sum(rate(container_cpu_usage_seconds_total{namespace="`) + namespace + string(`",container="`) + name + string(`"}[2m]))`)
	metrics := h.Auto.PromClient.GetPrometheus(string(queryExp))
	return metrics[0].Float()
}
func (h *Horization) ComputePodCpuLimit(namespace, name string) float64 {
	queryExp := string(`rate(kube_pod_container_resource_requests_cpu_cores{namespace="`) + namespace + string(`",container="`) + name + string(`"}[2m])`)
	metrics := h.Auto.PromClient.GetPrometheus(string(queryExp))
	return metrics[0].Float()
}

func (h *Horization) ComputeReplicasForMemoryUtilization(spec k8sv1beta1.K8sScaleSpec) (k8sv1beta1.K8sScaleSpec, bool) {
	var limit float64 = h.ComputeMemoryLimit(spec.Namespace, spec.Name)
	var usage float64 = h.ComputeMemoryUsage(spec.Namespace, spec.Name)
	//var replicas int32
	var tag bool
	spec.CurrentPercentage = usage / limit
	spec.CurrentReplicas, _ = h.Auto.K8sClient.K8sGetReplicas(spec.Namespace, spec.Name)
	spec.Replicas, _ = h.Auto.K8sClient.K8sGetReplicas(spec.Namespace, spec.Name)
	if spec.CurrentPercentage > spec.TargetPercentage {
		podlimit := spec.CurrentPercentage - spec.TargetPercentage
		podcurrentlimit := h.ComputePodMemoryLimit(spec.Namespace, spec.Name)
		replicas := podlimit / podcurrentlimit
		spec.Replicas = int32(math.Ceil(replicas))
		tag = true
		return spec, tag
	} else {
		tag = false
		return spec, tag
	}
}
func (h *Horization) ComputeMemoryLimit(namespace string, name string) float64 {
	queryExp := string(`sum(kube_pod_container_resource_requests_memory_cores{namespace="`) + namespace + string(`",container="`) + name + string(`"})`)
	metrics := h.Auto.PromClient.GetPrometheus(string(queryExp))
	return metrics[0].Float()
}
func (h *Horization) ComputeMemoryUsage(namespace string, name string) float64 {
	queryExp := string(`sum(container_memory_usage_bytes{namespace="`) + namespace + string(`",container="`) + name + string(`"})`)
	metrics := h.Auto.PromClient.GetPrometheus(string(queryExp))
	return metrics[0].Float()
}
func (h *Horization) ComputePodMemoryLimit(namespace, name string) float64 {
	queryExp := string(`kube_pod_container_resource_requests_memory_cores{namespace="`) + namespace + string(`",container="`) + name + string(`"}`)
	metrics := h.Auto.PromClient.GetPrometheus(string(queryExp))
	return metrics[0].Float()
}
func (h *Horization) ComputeReplicasForCustomMetrics() {}

func (h *Horization) Shouldscaler()          {}
func (h *Horization) UpdateCurrentReplicas() {}

func (h *Horization) UpdateStatus() {}
