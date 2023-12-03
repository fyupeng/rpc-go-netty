package client

import (
	"github.com/fyupeng/rpc-go-netty/config"
	"github.com/fyupeng/rpc-go-netty/discovery/load_balancer"
	"github.com/fyupeng/rpc-go-netty/discovery/service_discovery"
	"github.com/fyupeng/rpc-go-netty/factory"
	"github.com/fyupeng/rpc-go-netty/net/netty/future"
	"github.com/fyupeng/rpc-go-netty/protocol"
	"github.com/fyupeng/rpc-go-netty/utils/idworker"
	"log"
	"net/netip"
	"reflect"
	"strings"
	"time"
)

/*
* New NettyClient to Alone
 */
func NewNettyClientDirect(serviceAddress string, serializerCode int) RpcClient {
	worker, err := idworker.NewNowWorker(0)
	if err != nil {
		log.Fatal("id create failed: ", err)
	}
	return &nettyClient{
		ServiceConsumer: service_discovery.NewServiceConsumerDirect(serializerCode, serviceAddress),
		IdWorker:        worker,
		unProcessResult: (factory.GetInstance(reflect.TypeOf((*future.UnProcessResult)(nil)))).(*future.UnProcessResult),
		ServiceAddress:  serviceAddress,
	}
}

/*
* New NettyClient to Alone
 */
func NewNettyClient2Alone(loadBalancer load_balancer.LoadBalancer, serializerCode int, registryCenterAddress string) RpcClient {
	worker, err := idworker.NewNowWorker(0)
	if err != nil {
		log.Fatal("id create failed: ", err)
	}
	return &nettyClient{
		ServiceConsumer: service_discovery.NewServiceConsumerAlone(loadBalancer, serializerCode, registryCenterAddress),
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
	ServiceAddress  string
}

func (nettyClient *nettyClient) SendRequest(rpcRequest protocol.RequestProtocol) *future.CompleteFuture {
	groupName := rpcRequest.GetGroup()
	var serviceAddr netip.AddrPort
	var getServiceErr error

	serviceAddrString := strings.TrimSpace(nettyClient.ServiceAddress)
	if serviceAddrString != "" {
		serviceAddr, getServiceErr = config.ParseAddress(serviceAddrString)
	} else if groupName != "" {
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
