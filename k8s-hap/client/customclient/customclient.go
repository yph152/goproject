package customclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/k8s-hap/k8sv1beta1"
	"net/http"
)

type K8sResource struct {
	Namespace k8sv1beta1.K8sScaleSpec
}
type K8sSpec k8sv1beta1.K8sScaleSpec

type CustomClient struct {
	Client    *http.Client
	Endpoints string
}

type CustomClientInterface interface {
	Path(namespace string, name string) string
	Url(path string) string
	CustomGet(url string) K8sSpec
	CustomPost(url string, jsonstr []byte)
	CustomPut(url string, jsonstr []byte)
	CustomDelete(url string) K8sSpec
}

func (c *CustomClient) Path(namespace string, name string) string {
	if len(name) == 0 && len(namespace) == 0 {
		path := "/api/v1/namespaces"
		return path
	} else if len(name) == 0 {
		path := "/api/v1/autoscalers/" + namespace
		return path
	} else {
		path := "/api/v1/autoscalers/" + namespace + "/" + name
		return path
	}
}
func (c *CustomClient) Url(path string) string {
	return c.Endpoints + path
}

func (c *CustomClient) CustomGet(url string) k8sv1beta1.K8sScaleSpec {
	resp, err := http.Get(url)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("panic body\n")
	}

	var data K8sResource
	json.Unmarshal(body, &(data.Namespace))
	fmt.Println(data.Namespace)

	return data.Namespace
}
func (c *CustomClient) CustomGetAutoScalers(url string) []string {
	resp, err := http.Get(url)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("panic body\n")
	}
	var data []string
	json.Unmarshal(body, &(data))
	return data
}
func (c *CustomClient) CustomGetNameSpace(url string) []string {
	resp, err := http.Get(url)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("panic body\n")
	}
	var data []string
	json.Unmarshal(body, &(data))
	return data
}
func (c *CustomClient) CustomPost(url string, jsonstr []byte) {
	//	url := "http://127.0.0.1:8080/api/v1/namespaces/{namespace}/autoscalers/{name}"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonstr))
	req.Header.Set("Content-Type", "application/json")

	//	client := &http.Client{}
	resp, err := c.Client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

func (c *CustomClient) CustomPut(url string, jsonstr []byte) {

	//	url := "http://127.0.0.1:8080/api/v1/namespaces/{namespace}/autoscalers"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonstr))
	req.Header.Set("Content-Type", "application")

	//	client := &http.Client{}
	resp, err := c.Client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func (c *CustomClient) CustomDelete(url string) {
	//	usl := "http://127.0.0.1:8080/api/v1/namespaces/{namespace}/autoscalers/{name}"
	req, err := http.NewRequest("DELETE", url, nil)

	//	client := &http.Client{}
	resp, err := c.Client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
