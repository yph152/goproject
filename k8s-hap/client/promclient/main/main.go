package main

import (
	"fmt"
	"k8s.io/k8s-hap/client/promclient"
)

func main() {
	prom := promclient.NewPrometheusClient("10.3.240.8:9090")
	queryExp := `sum(rate(container_cpu_usage_seconds_total{namespace="gslb",container_name="jetty"}[2m]))`
	result := prom.GetPrometheus(string(queryExp))
	for _, value := range result {
		fmt.Println(value.Float())

	}
	fmt.Println("vim-go")
}
