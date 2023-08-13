package client

import (
	"rpc-go-netty/discovery/load_balancer"
	"rpc-go-netty/discovery/service_discovery"
)

/*
* New NettyClient to Alone
 */
func NewNettyClient2Alone(loadBalancer load_balancer.LoadBalancer, serializerCode int, registryCenterAddress string) service_discovery.ServiceDiscovery {
	consumer := service_discovery.NewServiceConsumer(loadBalancer, serializerCode, registryCenterAddress)
	return consumer
}

func NewNettyClient2Cluster(loadBalancer load_balancer.LoadBalancer, serializerCode int, registryCenterAddress string) service_discovery.ServiceDiscovery {
	consumer := service_discovery.NewServiceConsumer(loadBalancer, serializerCode, registryCenterAddress)
	return consumer
}

type client struct {
}
