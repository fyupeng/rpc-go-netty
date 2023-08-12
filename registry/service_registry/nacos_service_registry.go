package service_registry

import (
	"log"
	"reflect"
	"rpc-go-netty/codec"
	"rpc-go-netty/config"
	"rpc-go-netty/net/handler"
	"sync"
)

/*
*
服务提供者 实现 服务注册器接口（服务提供者拥有了服务注册的行为）
*/
func NewServiceProvider(serviceAddress, registerAddress string, serializerCode int) ServiceRegistry {

	services := make(map[string]interface{})

	serverHandler := handler.NewServerHandler(services)

	return &serviceProvider{
		ServerConfig:       config.NewServerConfig(serviceAddress, registerAddress, serverHandler, codec.CommonCodec(0, 8, serializerCode)),
		Services:           services,
		RegisteredServices: make(map[string]bool),
	}
}

type serviceProvider struct {
	ServerConfig       config.Config
	Services           map[string]interface{}
	RegisteredServices map[string]bool
	mutex              sync.Mutex
}

func (serviceProvider *serviceProvider) AddService(service interface{}, serviceName string) {
	serviceProvider.mutex.Lock()

	defer serviceProvider.mutex.Unlock()

	if _, ok := serviceProvider.RegisteredServices[serviceName]; ok {
		return
	}

	serviceProvider.RegisteredServices[serviceName] = true

	serviceProvider.Services[serviceName] = service

	serviceValue := reflect.ValueOf(service)

	log.Printf("Register cn.fyupeng.service: %v with interface %v", serviceValue, serviceName)

}

func (serviceProvider *serviceProvider) Listen() (err error) {
	err = serviceProvider.ServerConfig.Listen()
	return
}

func (serviceProvider *serviceProvider) Register(serviceName string) (err error) {
	err = serviceProvider.ServerConfig.RegisterInstance(serviceName)
	return
}

func (serviceProvider *serviceProvider) RegisterWithGroupName(serviceName string, groupName string) (err error) {
	err = serviceProvider.ServerConfig.RegisterInstanceWithGroupName(serviceName, groupName)
	return
}

func (serviceProvider *serviceProvider) ClearRegistry() (err error) {
	err = serviceProvider.ServerConfig.DeregisterAllInstance()
	return
}
