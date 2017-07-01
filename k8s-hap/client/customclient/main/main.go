package main

import (
	//	"fmt"
	"k8s.io/k8s-hap/client/customclient"
	"net/http"
)

func main() {
	c := customclient.CustomClient{Client: &http.Client{}, Endpoints: "http://127.0.0.1:8088"}

	path := c.Path("default", "svc")
	url := c.Url(path)
	//	c.CustomGet(url)

	//	c.CustomGet(url)
	str := `{"Name":"svc","Namespace":"default","Replicas":3}`
	str1 := `{"Name":"svc1","Namespace":"default","Replicas":2}`
	c.CustomPost(url, []byte(str))

	path = c.Path("default", "svc1")
	url = c.Url(path)
	c.CustomPost(url, []byte(str1))

	path = c.Path("default", "svc1")
	url = c.Url(path)
	c.CustomGet(url)
	//	c.CustomDelete(url)
	//	c.CustomGet(url)
	//	fmt.Println("vim-go")
}
