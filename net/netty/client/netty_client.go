package client

import (
	"log"
	"rpc-go-netty/discovery/load_balancer"
	"rpc-go-netty/discovery/service_discovery"
	"rpc-go-netty/protocol"
	"rpc-go-netty/utils/idworker"
	"strconv"
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
	worker, err := idworker.NewNowWorker(0)
	if err != nil {
		log.Fatal("id create failed: ", err)
	}
	return &nettyClient{
		ServiceConsumer: service_discovery.NewServiceConsumerWithCluster(loadBalancer, serializerCode, registryCenterAddress),
		IdWorker:        worker,
	}
}

type nettyClient struct {
	ServiceConsumer service_discovery.ServiceDiscovery
	IdWorker        idworker.IdWorker
}

func (nettyClient *nettyClient) SendRequest(interfaceName string, methodName string, parameters []interface{}, paramTypes []string, returnTypes []string) {
	serviceAddr, getServiceErr := nettyClient.ServiceConsumer.LookupServiceWithGroupName(interfaceName, "1.0.1")

	if getServiceErr != nil {
		log.Fatal("get Service Fatal: ", getServiceErr)
	}

	channel := nettyClient.ServiceConsumer.GetChannel(serviceAddr.String())

	message := protocol.RpcRequestProtocol(strconv.FormatInt(nettyClient.IdWorker.NextId(), 10), interfaceName, methodName, parameters,
		paramTypes, returnTypes[0], false, "1.0.1", false)

	err := channel.Write(message)
	if err != nil {
		log.Fatal("channel1 err: ", err)
	}
}
