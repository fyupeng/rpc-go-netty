package service_discovery

import (
	"github.com/go-netty/go-netty"
	"net/netip"
)

/*
*
服务发现接口
*/
type ServiceDiscovery interface {
	GetChannel(address string) netty.Channel

	LookupService(serviceName string) (netip.AddrPort, error)

	LookupServiceWithGroupName(serviceName string, groupName string) (netip.AddrPort, error)
}
