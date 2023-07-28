package service_registry

import (
	"github.com/go-netty/go-netty"
	"rpc-go-netty/config"
)

/*
*
服务提供者 实现 服务注册器接口（服务提供者拥有了服务注册的行为）
*/
func NewServiceProvider(serviceAddress, registerAddress string, serverHandler netty.ChannelHandler, commonCodec netty.CodecHandler) ServiceRegistry {

	return &serviceProvider{
		ServerConfig: config.NewServerConfig(serviceAddress, registerAddress, serverHandler, commonCodec),
	}
}

type serviceProvider struct {
	ServerConfig config.Config
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
