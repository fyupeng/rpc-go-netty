package client

import (
	"log"
	"net/netip"
	"reflect"
	"rpc-go-netty/discovery/load_balancer"
	"rpc-go-netty/discovery/service_discovery"
	"rpc-go-netty/factory"
	"rpc-go-netty/net/netty/future"
	"rpc-go-netty/protocol"
	"rpc-go-netty/utils/idworker"
	"time"
)

/*
* New NettyClient to Alone
 */
func NewNettyClient2Alone(loadBalancer load_balancer.LoadBalancer, serializerCode int, registryCenterAddress string) RpcClient {
	worker, err := idworker.NewNowWorker(0)
	if err != nil {
		log.Fatal("id create failed: ", err)
	}
	return &nettyClient{
		ServiceConsumer: service_discovery.NewServiceConsumer(loadBalancer, serializerCode, registryCenterAddress),
		IdWorker:        worker,
		unProcessResult: (factory.GetInstance(reflect.TypeOf((*future.UnProcessResult)(nil)))).(*future.UnProcessResult),
	}
}

func NewNettyClient2Cluster(loadBalancer load_balancer.LoadBalancer, serializerCode int, registryCenterAddress []string) RpcClient {
	return &nettyClient{
		ServiceConsumer: service_discovery.NewServiceConsumerWithCluster(loadBalancer, serializerCode, registryCenterAddress),
		unProcessResult: (factory.GetInstance(reflect.TypeOf((*future.UnProcessResult)(nil)))).(*future.UnProcessResult),
	}
}

type nettyClient struct {
	ServiceConsumer service_discovery.ServiceDiscovery
	unProcessResult *future.UnProcessResult
	IdWorker        idworker.IdWorker
}

func (nettyClient *nettyClient) SendRequest(rpcRequest protocol.RequestProtocol) *future.CompleteFuture {
	groupName := rpcRequest.GetGroup()
	var serviceAddr netip.AddrPort
	var getServiceErr error
	if groupName != "" {
		serviceAddr, getServiceErr = nettyClient.ServiceConsumer.LookupServiceWithGroupName(rpcRequest.GetInterfaceName(), groupName)
	} else {
		serviceAddr, getServiceErr = nettyClient.ServiceConsumer.LookupService(rpcRequest.GetInterfaceName())
	}

	if getServiceErr != nil {
		log.Fatal("get Service Fatal: ", getServiceErr)
	}

	channel := nettyClient.ServiceConsumer.GetChannel(serviceAddr.String())

	message := protocol.RpcRequestProtocol(rpcRequest.GetRequestId(), rpcRequest.GetInterfaceName(), rpcRequest.GetMethodName(), rpcRequest.GetParameters(),
		rpcRequest.GetParamTypes(), rpcRequest.GetReturnType(), false, rpcRequest.GetGroup(), false)

	completeFuture := future.NewCompleteFuture(make(chan interface{}), time.Second*10)
	//unProcessResult := NewUnprocessResult()
	nettyClient.unProcessResult.Put(rpcRequest.GetRequestId(), completeFuture)

	err := channel.Write(message)
	if err != nil {
		log.Fatal("channel1 err: ", err)
	}

	return completeFuture

}
