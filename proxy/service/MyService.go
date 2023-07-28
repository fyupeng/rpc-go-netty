package service

import "fmt"

type ServiceAnnotation struct {
	// 定义你的注解字段
	// ...
}

type Service interface {
	Hello()
}

// @Register("service", "arg1", "arg2")
func MyService() Service {
	return &myService{
		ServiceAnnotation: ServiceAnnotation{},
	}
}

type myService struct {
	ServiceAnnotation
	// 其他字段...
}

func (myService *myService) Hello() {
	fmt.Println("hello你好")
}
