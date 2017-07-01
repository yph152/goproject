package k8sapiserver

import (
	//"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/pkg/api"
	//"k8s.io/client-go/pkg/api/unversioned"
	//	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
	//	"k8s.io/k8s-hap/client/customclient"
	//"k8s.io/k8s-hap/k8sv1beta1"
	//	"log"
)

/*var (
	kubeconfig = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
)*/

type K8sapiserver struct {
	K8sapi *kubernetes.Clientset
}

func NewClient(kubeconfig string) *K8sapiserver {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	k8sapiserver := K8sapiserver{K8sapi: clientset}
	return &k8sapiserver
}
func (k8s *K8sapiserver) K8sGetReplicas(namespace, name string) (int32, error) {
	deployments, err := k8s.K8sapi.Extensions().Deployments(namespace).Get(name)
	spec_num := deployments.Spec.Replicas
	return *spec_num, err
}
func (k8s *K8sapiserver) K8sScaler(namespace, name string, scalereplicas int32) error {
	scale, err := k8s.K8sapi.Extensions().Scales(namespace).Get("deployment", name)

	check_err(err)
	scale.Spec.Replicas = scalereplicas
	_, err = k8s.K8sapi.Extensions().Scales(namespace).Update("deployment", scale)
	check_err(err)
	return err
}

func check_err(err error) {
	if err != nil {
		fmt.Printf("go err from apiserver: %s\n", err)
	}
}
