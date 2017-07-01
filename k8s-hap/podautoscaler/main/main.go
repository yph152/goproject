package main

import (
	"fmt"
	"k8s.io/k8s-hap/podautoscaler/queue"
)

type Test1 struct {
	Apple int
	Bana  string
}

func (test *Test1) get_a() {
	fmt.Printf("%d\n", test.Apple)
}

func (test *Test1) get_b() {
	fmt.Printf("%s\n", test.Bana)
}

func worker(val interface{}) {
	switch v := val.(type) {
	case Test1:
		v.get_a()
		v.get_b()
	default:
		println("no have this type")
	}
	//	fmt.Println(val.b)
}
func main() {
	q := queue.NewQueue(worker, 1)
	s := Test1{}
	for i := 0; i < 10; i++ {
		s.Apple = i
		s.Bana = "hello"
		q.Push(s)
	}
	q.Wait()
	fmt.Println("vim-go")
}
