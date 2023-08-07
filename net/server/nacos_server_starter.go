package server

import (
	"log"
	"rpc-go-netty/codec"
	"rpc-go-netty/registry/service_registry"
)

func NewNacosServerStarter(serviceAddress, registerAddress string) ServerStart {

	return &nacosServerStarter{
		serverAddress:   serviceAddress,
		serviceRegistry: service_registry.NewServiceProvider(serviceAddress, registerAddress, NewServerHandler(), codec.CommonCodec(0, 1024, 0)),
	}

}

type nacosServerStarter struct {
	serverAddress   string
	serviceRegistry service_registry.ServiceRegistry
}

func (server *nacosServerStarter) Start() (err error) {
	// 监听端口
	server.publishService()
	err = server.serviceRegistry.Listen()
	return
}

func (server *nacosServerStarter) publishService() {
	err := server.serviceRegistry.RegisterWithGroupName("cn.fyupeng.service.HelloService", "1.0.0")
	if err != nil {
		log.Fatal("publish fatail: ", err)
	}
	return
}
