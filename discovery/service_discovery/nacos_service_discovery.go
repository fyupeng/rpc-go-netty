package service_discovery

import (
	"github.com/fyupeng/rpc-go-netty/codec"
	"github.com/fyupeng/rpc-go-netty/config"
	"github.com/fyupeng/rpc-go-netty/discovery/load_balancer"
	"github.com/fyupeng/rpc-go-netty/net/handler"
	"github.com/go-netty/go-netty"
	"github.com/nacos-group/nacos-sdk-go/model"
	"net"
	"net/netip"
	"strconv"
)

/*
*
服务消费者 实现 服务发现接口（服务消费者拥有了服务发现的行为）
*/
func NewServiceConsumerDirect(serializerCode int, serviceAddress string) ServiceDiscovery {

	return &serviceConsumer{
		ClientConfig:   config.NewClientConfigDirect(handler.NewClientHandler(), codec.CommonCodec(0, 8, serializerCode), handler.NewRequestParser(serializerCode)),
		ServiceAddress: serviceAddress,
	}
}

func NewServiceConsumerAlone(loadBalancer load_balancer.LoadBalancer, serializerCode int, registryCenterAddress string) ServiceDiscovery {
	var registryCenterAddressArray []string
	registryCenterAddressArray = append(registryCenterAddressArray, registryCenterAddress)

	return &serviceConsumer{
		ClientConfig: config.NewClientConfig(registryCenterAddressArray, handler.NewClientHandler(), codec.CommonCodec(0, 8, serializerCode), handler.NewRequestParser(serializerCode)),
		LoadBalancer: loadBalancer,
	}
}

func NewServiceConsumerWithCluster(loadBalancer load_balancer.LoadBalancer, serializerCode int, registryCenterAddress []string) ServiceDiscovery {

	return &serviceConsumer{
		ClientConfig: config.NewClientConfig(registryCenterAddress, handler.NewClientHandler(), codec.CommonCodec(0, 8, serializerCode), handler.NewRequestParser(serializerCode)),
		LoadBalancer: loadBalancer,
	}
}

type serviceConsumer struct {
	ClientConfig   config.Config
	LoadBalancer   load_balancer.LoadBalancer
	ServiceAddress string
}

func (serviceConsumer *serviceConsumer) GetChannel(address string) netty.Channel {
	return serviceConsumer.ClientConfig.GetChannel(address)
}

func (serviceConsumer *serviceConsumer) LookupService(serviceName string) (netip.AddrPort, error) {
	services, err := serviceConsumer.ClientConfig.GetAllInstance(serviceName)

	if err != nil {
		return netip.AddrPort{}, err
	}

	service, selectServiceErr := serviceConsumer.LoadBalancer.SelectService(services)
	if selectServiceErr != nil {
		return netip.AddrPort{}, err
	}

	addrPort, err := config.ParseInstance2AddPort((service).(model.Instance))

	return addrPort, err

}

func (serviceConsumer *serviceConsumer) LookupServiceWithGroupName(serviceName string, groupName string) (netip.AddrPort, error) {
	services, err := serviceConsumer.ClientConfig.GetAllInstanceWithGroupName(serviceName, groupName)

	if err != nil {
		return netip.AddrPort{}, err
	}

	service, selectServiceErr := serviceConsumer.LoadBalancer.SelectService(services)
	if selectServiceErr != nil {
		return netip.AddrPort{}, err
	}

	instance := (service).(model.Instance)

	ip, _ := net.ResolveIPAddr("ip", instance.Ip)
	port := strconv.FormatUint(instance.Port, 10)

	addrPort, _ := netip.ParseAddrPort(ip.String() + ":" + port)

	return addrPort, nil
}
