//usr/bin/env go run $0 "$@"; exit
package promclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	//"os"

	//apiv1 "k8s.io/client-go/pkg/api/v1"

	//"github.com/golang/glog"
	"github.com/tidwall/gjson"
	//	"github.com/zanhsieh/custom-k8s-hpa/k8sclient"
)

type PrometheusClient struct {
	prometheus_ip_port string
}

func NewPrometheusClient(endpoints string) *PrometheusClient {
	return &PrometheusClient{
		prometheus_ip_port: endpoints,
	}
}

func (s *PrometheusClient) GetPrometheus(queryExp string) []gjson.Result {
	serverPath := "/api/v1/query?query=%v"
	tmp := fmt.Sprintf(serverPath, queryExp)
	url := fmt.Sprintf("http://%v%v", s.prometheus_ip_port, tmp)

	res, err := http.Get(url)

	if err != nil {
		fmt.Println("Get error")
		return nil
	}

	jsonString, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println("ReadAll error")
		return nil
	}

	pre_metric := gjson.Get(string(jsonString), "data.result.#.value")
	if len(pre_metric.Array()) == 0 {
		return nil
	}
	metric := gjson.Get(string(jsonString), "data.result.#.value.1")

	return metric.Array()
}
