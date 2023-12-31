package server

import (
	"github.com/fyupeng/rpc-go-netty/net/handler"
	"github.com/fyupeng/rpc-go-netty/provider/service_provider"
	"github.com/fyupeng/rpc-go-netty/registry/service_registry"
	"log"
	"reflect"
)

func NewNettyServer(serviceAddress, registerAddress string, serializerCode int) RpcServer {
	var registerAddressArray []string
	registerAddressArray = append(registerAddressArray, registerAddress)

	serviceProvider := service_provider.NewDefaultServiceProvider()

	serverHandler := handler.NewServerHandler(serviceProvider, serializerCode)

	serviceRegistry := service_registry.NewNacosServiceRegistry(serviceAddress, registerAddressArray, serverHandler, serializerCode)

	return &nettyServerStarter{
		serverAddress:   serviceAddress,
		serviceRegistry: serviceRegistry,
		serviceProvider: serviceProvider,
	}

}

func NewNettyServerWithCluster(serviceAddress string, registerAddress []string, serializerCode int) RpcServer {

	serviceProvider := service_provider.NewDefaultServiceProvider()

	serverHandler := handler.NewServerHandler(serviceProvider, serializerCode)

	serviceRegistry := service_registry.NewNacosServiceRegistry(serviceAddress, registerAddress, serverHandler, serializerCode)

	return &nettyServerStarter{
		serverAddress:   serviceAddress,
		serviceRegistry: serviceRegistry,
		serviceProvider: serviceProvider,
	}

}

type nettyServerStarter struct {
	serverAddress   string
	serviceRegistry service_registry.ServiceRegistry
	serviceProvider service_provider.ServiceProvider
}

func (server *nettyServerStarter) Start() (err error) {
	// 监听端口
	err = server.serviceRegistry.Listen()
	return
}

func (server *nettyServerStarter) PublishService(services ...interface{}) {
	for _, service := range services {
		structType := reflect.TypeOf(service).Elem()

		var groupName string

		interfaceName := structType.Name()

		if nameField, isValid := structType.FieldByName("name"); isValid {
			interfaceName = nameField.Tag.Get("annotation")
			log.Printf("Found service %s which Field: %s, Annotation: %s\n", structType, "name", interfaceName)
		}

		if groupField, isValid := structType.FieldByName("group"); isValid {
			groupName = groupField.Tag.Get("annotation")
			log.Printf("Found service %s which Field: %s, Annotation: %s\n", structType, "group", groupName)
		}

		server.serviceProvider.AddService(service, interfaceName)

		err := server.serviceRegistry.RegisterWithGroupName(interfaceName, groupName)
		if err != nil {
			log.Fatal("interfaceName publish failed: ", err)
		}

		log.Printf("Register service %v with interface %v\n", structType, interfaceName)

	}

	return
}
