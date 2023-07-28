package service_discovery

import (
	"github.com/go-netty/go-netty"
	"github.com/nacos-group/nacos-sdk-go/model"
	"net"
	"net/netip"
	"rpc-go-netty/config"
	"rpc-go-netty/discovery/load_balancer"
	"strconv"
)

/*
*
服务消费者 实现 服务发现接口（服务消费者拥有了服务发现的行为）
*/
func NewServiceConsumer(loadBalancer load_balancer.LoadBalancer, registryCenterAddress string, clientHandler netty.ChannelHandler, commonCodec netty.CodecHandler) ServiceDiscovery {

	return &serviceConsumer{
		ClientConfig: config.NewClientConfig(registryCenterAddress, clientHandler, commonCodec),
		LoadBalancer: loadBalancer,
	}
}

type serviceConsumer struct {
	ClientConfig config.Config
	LoadBalancer load_balancer.LoadBalancer
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
