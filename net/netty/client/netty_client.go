package client

import (
	"log"
	"rpc-go-netty/discovery/load_balancer"
	"rpc-go-netty/discovery/service_discovery"
	"rpc-go-netty/protocol"
)

/*
* New NettyClient to Alone
 */
func NewNettyClient2Alone(loadBalancer load_balancer.LoadBalancer, serializerCode int, registryCenterAddress string) RpcClient {
	return &nettyClient{
		ServiceConsumer: service_discovery.NewServiceConsumer(loadBalancer, serializerCode, registryCenterAddress),
	}
}

func NewNettyClient2Cluster(loadBalancer load_balancer.LoadBalancer, serializerCode int, registryCenterAddress []string) RpcClient {
	return &nettyClient{
		ServiceConsumer: service_discovery.NewServiceConsumerWithCluster(loadBalancer, serializerCode, registryCenterAddress),
	}
}

type nettyClient struct {
	ServiceConsumer service_discovery.ServiceDiscovery
}

func (nettyClient *nettyClient) SendRequest(interfaceName string, methodName string, parameters []interface{}, paramTypes []string, returnTypes []string) {
	serviceAddr, getServiceErr := nettyClient.ServiceConsumer.LookupServiceWithGroupName(interfaceName, "1.0.1")

	if getServiceErr != nil {
		log.Fatal("get Service Fatal: ", getServiceErr)
	}

	channel := nettyClient.ServiceConsumer.GetChannel(serviceAddr.String())

	message := protocol.RpcRequestProtocol("123455", interfaceName, methodName, parameters,
		paramTypes, returnTypes[0], false, "1.0.1", false)

	err := channel.Write(message)
	if err != nil {
		log.Fatal("channel1 err: ", err)
	}
}
