package config

import (
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec/frame"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

type Config interface {
	GetChannel(address string) netty.Channel

	Listen() error

	GetAllInstance(serviceName string) ([]model.Instance, error)

	GetAllInstanceWithGroupName(serviceName string, groupName string) ([]model.Instance, error)

	RegisterInstance(serviceName string) (err error)

	RegisterInstanceWithGroupName(serviceName string, groupName string) (err error)

	DeregisterAllInstance() (err error)

	DeregisterAllInstanceWithGroupName(groupName string) (err error)
}

func NewServerConfig(serviceAddress, registryServerAddress string, serverHandler netty.ChannelHandler, commonCodec netty.CodecHandler) Config {

	serviceAddrPort, err := ParseAddress(serviceAddress)
	registryServerAddrPort, err := ParseAddress(registryServerAddress)

	if err != nil {
		log.Fatal("new client config fatal: ", err)
	}

	clientConfig := createClientConfig()
	serverConfigs := createServerConfig(registryServerAddrPort)
	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	return &serviceConfig{
		RegistryServerAddrPort: registryServerAddrPort,
		ServiceAddrPort:        serviceAddrPort,
		ServiceNames:           make(map[string]bool),
		ServerHandler:          serverHandler,
		CommonCodec:            commonCodec,
		NamingClient:           namingClient,
	}
}

// 创建RPC客户端和服务端服务配置
type serviceConfig struct {
	// 注册中心服务地址
	RegistryServerAddrPort netip.AddrPort
	// 服务提供地址
	ServiceAddrPort netip.AddrPort
	// 服务名
	ServiceNames map[string]bool
	// RPC服务端/客户端 client
	NamingClient naming_client.INamingClient
	// 客户端 连接服务端 通道
	Channels map[string]netty.Channel
	// 编码器
	CommonCodec netty.CodecHandler
	// 服务端处理器
	ServerHandler netty.ChannelHandler
	// 客户端处理器
	ClientHandler netty.ChannelHandler
	// 客户端 引导
	ClientBootstrap netty.Bootstrap
}

func NewClientConfig(registryCenterAddress string, clientHandler netty.ChannelHandler, commonCodec netty.CodecHandler, protocolHandler netty.ChannelHandler) Config {

	address, err := ParseAddress(registryCenterAddress)

	if err != nil {
		log.Fatal("new client config fatal: ", err)
	}

	clientConfig := createClientConfig()
	serverConfigs := createServerConfig(address)
	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	return &serviceConfig{
		NamingClient:           namingClient,
		RegistryServerAddrPort: address,
		ClientHandler:          clientHandler,
		CommonCodec:            commonCodec,
		Channels:               make(map[string]netty.Channel),
		ClientBootstrap:        initClientBootstrap(clientHandler, commonCodec, protocolHandler),
	}

}

func createClientConfig() constant.ClientConfig {
	// 创建clientConfig的另一种方式
	clientConfig := constant.NewClientConfig(
		constant.WithNamespaceId(""), //当namespace是public时，此处填空字符串。
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)
	return *clientConfig
}

func createServerConfig(address netip.AddrPort) []constant.ServerConfig {
	// 创建serverConfig的另一种方式
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			address.Addr().String(),
			uint64(address.Port()),
			constant.WithScheme("http"),
			constant.WithContextPath("/nacos"),
		),
	}
	return serverConfigs
}

func (serviceConfig *serviceConfig) Listen() (err error) {
	// 监听端口
	bootstrap := initServerBootstrap(serviceConfig.ServerHandler, serviceConfig.CommonCodec)
	listener := bootstrap.Listen(serviceConfig.ServiceAddrPort.String())
	log.Println("listening port... [" + serviceConfig.ServiceAddrPort.String() + "]")
	err = listener.Sync()
	return
}

func initServerBootstrap(serverHandler netty.ChannelHandler, commonCodec netty.CodecHandler) netty.Bootstrap {
	var childInitializer = func(channel netty.Channel) {
		channel.Pipeline().
			//AddLast(netty.ReadIdleHandler(time.Second * 3)).
			//AddLast(netty.WriteIdleHandler(time.Second * 5)).
			AddLast(frame.DelimiterCodec(1024, "\r\n", true)).
			AddLast(commonCodec).
			AddLast(serverHandler)
	}

	return netty.NewBootstrap(netty.WithChildInitializer(childInitializer))

}

func initClientBootstrap(clientHandler netty.ChannelHandler, commonCodec netty.CodecHandler, protocolHandler netty.ChannelHandler) netty.Bootstrap {

	clientInitializer := func(channel netty.Channel) {
		channel.Pipeline().
			//AddLast(netty.ReadIdleHandler(time.Second * 3)).
			//AddLast(netty.WriteIdleHandler(time.Second * 5)).
			AddLast(frame.DelimiterCodec(1024, "\r\n", true)).
			AddLast(commonCodec).
			AddLast(protocolHandler).
			AddLast(clientHandler)

	}

	return netty.NewBootstrap(netty.WithClientInitializer(clientInitializer))

}

func (serviceConfig *serviceConfig) GetChannel(address string) netty.Channel {
	addrPort, err := ParseAddress(address)

	if err != nil {
		log.Fatal("parse fatail: ", err)
	}
	if serviceConfig.ClientBootstrap == nil {
		log.Fatal("please InitBootstrap first...")
	}

	channelKey := addrPort.String()

	channel := serviceConfig.Channels[channelKey]

	if channel != nil && channel.IsActive() {
		return channel
	}

	ch, err := serviceConfig.ClientBootstrap.Connect(addrPort.String())

	if err != nil {
		log.Fatal("serviceChannelProvider connect failed", err)
	}

	serviceConfig.Channels[channelKey] = ch

	return ch
}

func (serviceConfig *serviceConfig) GetAllInstance(serviceName string) ([]model.Instance, error) {
	instances, err := serviceConfig.NamingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		HealthyOnly: true,
	})
	return instances, err
}

func (serviceConfig *serviceConfig) GetAllInstanceWithGroupName(serviceName string, groupName string) ([]model.Instance, error) {
	instances, err := serviceConfig.NamingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   groupName,
		HealthyOnly: true,
	})
	return instances, err
}

func (serviceConfig *serviceConfig) RegisterInstance(serviceName string) (err error) {
	log.Println("正在注册服务【serviceName:%s, address:%s】", serviceName, serviceConfig.ServiceAddrPort.String())
	// 保存服务地址
	serviceConfig.ServiceAddrPort = serviceConfig.ServiceAddrPort
	// 保存服务名
	serviceConfig.ServiceNames[serviceName] = true
	_, err = serviceConfig.NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          serviceConfig.ServiceAddrPort.Addr().String(),
		Port:        uint64(serviceConfig.ServiceAddrPort.Port()),
		ServiceName: serviceName,
		Weight:      1,
		Healthy:     true,
		Enable:      true,
	})
	return
}

func (serviceConfig *serviceConfig) RegisterInstanceWithGroupName(serviceName string, groupName string) (err error) {
	// 保存服务地址
	serviceConfig.ServiceAddrPort = serviceConfig.ServiceAddrPort
	// 保存服务名
	serviceConfig.ServiceNames[serviceName] = true
	_, err = serviceConfig.NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          serviceConfig.ServiceAddrPort.Addr().String(),
		Port:        uint64(serviceConfig.ServiceAddrPort.Port()),
		ServiceName: serviceName,
		GroupName:   groupName,
		Weight:      1,
		Healthy:     true,
		Enable:      true,
	})
	return
}

func (serviceConfig *serviceConfig) DeregisterAllInstance() (err error) {
	for serviceName, _ := range serviceConfig.ServiceNames {
		_, err = serviceConfig.NamingClient.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          serviceConfig.ServiceAddrPort.Addr().String(),
			Port:        uint64(serviceConfig.ServiceAddrPort.Port()),
			ServiceName: serviceName,
		})
		delete(serviceConfig.ServiceNames, serviceName)
	}
	return
}

func (serviceConfig *serviceConfig) DeregisterAllInstanceWithGroupName(groupName string) (err error) {
	for serviceName, _ := range serviceConfig.ServiceNames {
		_, err = serviceConfig.NamingClient.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          serviceConfig.ServiceAddrPort.Addr().String(),
			Port:        uint64(serviceConfig.ServiceAddrPort.Port()),
			ServiceName: serviceName,
			GroupName:   groupName,
		})
	}
	return
}

func ParseAddress(registryCenterAddress string) (address netip.AddrPort, err error) {
	// 解析 域名
	addrPortArray := strings.Split(registryCenterAddress, ":")
	addr, _ := net.ResolveIPAddr("ip", addrPortArray[0])
	port := addrPortArray[1]
	address, err = netip.ParseAddrPort(addr.String() + ":" + port)
	return
}

func ParseInstance2AddPort(instance model.Instance) (addrPort netip.AddrPort, err error) {
	addrPort, err = netip.ParseAddrPort(instance.Ip + ":" + strconv.FormatUint(instance.Port, 10))
	return
}
