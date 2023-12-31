package service_registry

import (
	"github.com/fyupeng/rpc-go-netty/codec"
	"github.com/fyupeng/rpc-go-netty/config"
	"github.com/fyupeng/rpc-go-netty/net/handler"
	"github.com/go-netty/go-netty"
)

/*
*
服务提供者 实现 服务注册器接口（服务提供者拥有了服务注册的行为）
*/
func NewNacosServiceRegistry(serviceAddress string, registerAddress []string, serverHandler netty.ChannelHandler, serializerCode int) ServiceRegistry {

	return &nacosServiceRegistry{
		ServerConfig:       config.NewServerConfig(serviceAddress, registerAddress, serverHandler, codec.CommonCodec(0, 8, serializerCode)),
		RegisteredServices: make(map[string]bool),
		serializerCode:     serializerCode,
	}
}

type nacosServiceRegistry struct {
	ServerConfig       config.Config
	RegisteredServices map[string]bool
	serializerCode     int
}

func (serviceProvider *nacosServiceRegistry) Listen() (err error) {
	err = serviceProvider.ServerConfig.Listen(handler.NewResponseParser(serviceProvider.serializerCode))
	return
}

func (serviceProvider *nacosServiceRegistry) Register(serviceName string) (err error) {
	err = serviceProvider.ServerConfig.RegisterInstance(serviceName)
	return
}

func (serviceProvider *nacosServiceRegistry) RegisterWithGroupName(serviceName string, groupName string) (err error) {
	err = serviceProvider.ServerConfig.RegisterInstanceWithGroupName(serviceName, groupName)
	return
}

func (serviceProvider *nacosServiceRegistry) ClearRegistry() (err error) {
	err = serviceProvider.ServerConfig.DeregisterAllInstance()
	return
}
