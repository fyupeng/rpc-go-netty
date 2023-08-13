package service_provider

import (
	"log"
	"reflect"
	"sync"
)

/*
*
服务提供者 实现 服务注册器接口（服务提供者拥有了服务注册的行为）
*/
func NewDefaultServiceProvider() ServiceProvider {

	return &defaultServiceProvider{
		Services: make(map[string]interface{}),
	}
}

type defaultServiceProvider struct {
	Services map[string]interface{}
	mutex    sync.Mutex
}

func (serviceProvider *defaultServiceProvider) GetService(serviceName string) (service interface{}) {

	serviceProvider.mutex.Lock()

	defer serviceProvider.mutex.Unlock()

	if service, ok := serviceProvider.Services[serviceName]; !ok {
		log.Fatalf("Service[%v] Not Found!", serviceName)
	} else {
		return service
	}
	return
}

func (serviceProvider *defaultServiceProvider) AddService(service interface{}, serviceName string) {
	serviceProvider.mutex.Lock()

	defer serviceProvider.mutex.Unlock()

	if _, ok := serviceProvider.Services[serviceName]; ok {
		return
	}

	serviceProvider.Services[serviceName] = true

	serviceProvider.Services[serviceName] = service

	serviceValue := reflect.ValueOf(service)

	log.Printf("Register cn.fyupeng.service: %v with interface %v", serviceValue, serviceName)

}
