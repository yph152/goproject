package main

import (
	"flag"
	"fmt"
	//	"github.com/coreos/etcd/client"
	"k8s.io/k8s-hap/client/etcdclient"
)

var (
	etcdconfig = flag.String("etcdendpoints", "http://127.0.0.1:2379", "endpoints to the etcd")
)

func main() {
	flag.Parse()
	str1 := make([]string, 3)
	str1[0] = *etcdconfig
	client, err := etcdclient.NewClient(str1)
	c := etcdclient.EtcdClient{}
	c.Client = client
	//	str := `{"kind":"Namespace","apiVersion":"v1","metadata":{"name":"lyra","uid":"a20d22ae-a55b-11e6-9719-ecf4bbe850c8","creationTimestamp":"2016-11-08T02:32:46Z"},"spec":{"finalizers":["kubernetes"]},"status":{"phase":"Active"}}`
	//	fmt.Println(str)

	//	err = c.EtcdSet("/foo1/", str)
	/*	if err != nil {
			fmt.Print(err)
		}
	*/
	str, err := c.EtcdList("/foo")
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Println(str)
	}

}
