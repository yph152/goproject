//usr/bin/env go run $0 "$@"; exit
package options

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

type AutoscalerConfig struct {
	CustomClient     string
	Kubeconfig       string
	PrometheusIPPort string
}

func NewAutoScalerConfig() *AutoscalerConfig {
	return &AutoscalerConfig{
		CustomClient:     "127.0.0.1:8088",
		Kubeconfig:       "./config",
		PrometheusIPPort: "127.0.0.1:9090",
	}
}

func (c *AutoscalerConfig) ValidateFlags() error {
	var errorsFound bool
	if c.PrometheusIPPort == "" {
		errorsFound = true
		fmt.Println("prometheus cannot empty")
	}
	if c.CustomClient == "" {
		errorsFound = true
		glog.Errorf("Customapiserver  parameter cannot be empty")
	}
	if errorsFound {
		return fmt.Errorf("failed to validate all input parameters")
	}
	return nil
}

func (c *AutoscalerConfig) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.PrometheusIPPort, "prom-ip-port", c.PrometheusIPPort, "Prometheus ip/dns name and listen port number.")
	fs.StringVar(&c.CustomClient, "customapiserver", c.CustomClient, "Customapiserver ip and port .")
	fs.StringVar(&c.Kubeconfig, "kubeconfig", c.Kubeconfig, "K8sAPiserver used kubeconfig.")
}
