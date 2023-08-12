package server

import (
	"log"
	"reflect"
	"rpc-go-netty/registry/service_registry"
)

func NewNacosServerStarter(serviceAddress, registerAddress string, serializerCode int) RpcServer {

	serviceProvider := service_registry.NewServiceProvider(serviceAddress, registerAddress, serializerCode)

	return &nacosServerStarter{
		serverAddress:   serviceAddress,
		serviceRegistry: serviceProvider,
	}

}

type nacosServerStarter struct {
	serverAddress   string
	serviceRegistry service_registry.ServiceRegistry
}

func (server *nacosServerStarter) Start() (err error) {
	// 监听端口
	err = server.serviceRegistry.Listen()
	return
}

func (server *nacosServerStarter) PublishService(services ...interface{}) {
	for _, service := range services {
		structType := reflect.TypeOf(service).Elem()

		var groupName string

		interfaceName := structType.Name()

		if nameField, isValid := structType.FieldByName("name"); isValid {
			interfaceName = nameField.Tag.Get("annotation")
		}

		if groupField, isValid := structType.FieldByName("group"); isValid {
			groupName = groupField.Tag.Get("annotation")
		}

		server.serviceRegistry.AddService(service, interfaceName)

		err := server.serviceRegistry.RegisterWithGroupName(interfaceName, groupName)
		if err != nil {
			log.Fatal("publish fatail: ", err)
		}
	}

	return
}
