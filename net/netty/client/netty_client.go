package client

import (
	"fmt"
	"log"
	"reflect"
	"rpc-go-netty/discovery/load_balancer"
	"rpc-go-netty/discovery/service_discovery"
	"rpc-go-netty/factory"
	"rpc-go-netty/net/netty/future"
	"rpc-go-netty/protocol"
	"rpc-go-netty/utils/idworker"
	"strconv"
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
	worker, err := idworker.NewNowWorker(0)
	if err != nil {
		log.Fatal("id create failed: ", err)
	}
	return &nettyClient{
		ServiceConsumer: service_discovery.NewServiceConsumerWithCluster(loadBalancer, serializerCode, registryCenterAddress),
		IdWorker:        worker,
		unProcessResult: (factory.GetInstance(reflect.TypeOf((*future.UnProcessResult)(nil)))).(*future.UnProcessResult),
	}
}

type nettyClient struct {
	ServiceConsumer service_discovery.ServiceDiscovery
	IdWorker        idworker.IdWorker
	unProcessResult *future.UnProcessResult
}

func (nettyClient *nettyClient) SendRequest(interfaceName string, methodName string, parameters []interface{}, paramTypes []string, returnTypes []string) *future.CompleteFuture {
	serviceAddr, getServiceErr := nettyClient.ServiceConsumer.LookupServiceWithGroupName(interfaceName, "1.0.1")

	if getServiceErr != nil {
		log.Fatal("get Service Fatal: ", getServiceErr)
	}

	channel := nettyClient.ServiceConsumer.GetChannel(serviceAddr.String())

	requestId := strconv.FormatInt(nettyClient.IdWorker.NextId(), 10)

	message := protocol.RpcRequestProtocol(requestId, interfaceName, methodName, parameters,
		paramTypes, returnTypes[0], false, "1.0.1", false)

	completeFuture := future.NewCompleteFuture(make(chan interface{}), time.Second*10)
	//unProcessResult := NewUnprocessResult()
	fmt.Println("nettyClient.unProcessResult")
	fmt.Println(nettyClient.unProcessResult)
	nettyClient.unProcessResult.Put(requestId, completeFuture)

	err := channel.Write(message)
	if err != nil {
		log.Fatal("channel1 err: ", err)
	}

	return completeFuture

}
