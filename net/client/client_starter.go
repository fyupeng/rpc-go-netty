package client

import (
	"rpc-go-netty/discovery/load_balancer"
	"rpc-go-netty/discovery/service_discovery"
)

func NewClient(loadBalancer load_balancer.LoadBalancer, serializerCode int, registryCenterAddress string) service_discovery.ServiceDiscovery {
	consumer := service_discovery.NewServiceConsumer(loadBalancer, serializerCode, registryCenterAddress)
	return consumer
}

type client struct {
}
