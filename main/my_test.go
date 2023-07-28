package main

import (
	"fmt"
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty/codec/frame"
	"log"
	"rpc-go-netty/codec"
	"rpc-go-netty/discovery/load_balancer"
	"rpc-go-netty/discovery/service_discovery"
	"rpc-go-netty/net/client"
	"rpc-go-netty/net/server"
	"rpc-go-netty/protocol"
	"testing"
	"time"
)

func TestClient(t *testing.T) {

	serviceConsumer := service_discovery.NewServiceConsumer(load_balancer.LoadBalancer{}, "127.0.0.1:8848", client.NewClientHandler(), codec.CommonCodec(0, 1024, 71, 0))

	fmt.Println(serviceConsumer)

	//serviceAddr, getServiceErr := serviceConsumer.LookupService("TestService")
	//
	//if getServiceErr != nil {
	//	log.Fatal("get Service Fatal: ", getServiceErr)
	//}

	serviceAddr, getServiceErr := serviceConsumer.LookupServiceWithGroupName("helloService", "1.0.0")

	if getServiceErr != nil {
		log.Fatal("get Service Fatal: ", getServiceErr)
	}

	fmt.Println("serviceAddr is " + serviceAddr.String())

	channel := serviceConsumer.GetChannel("127.0.0.1:9527")

	fmt.Println("channel: ", channel)

	parameters := []interface{}{"hello，这里是go语言"}

	message := protocol.RpcRequestProtocol("123455", "helloService", "sayHello", parameters,
		[]string{"java.lang.String"}, "java.lang.String", false, "1.0.0", false)

	err := channel.Write(message)
	if err != nil {
		log.Fatal("channel1 err: ", err)
	}

	time.Sleep(time.Second * 60)

	// 在这里可以进行其他操作，发送消息等

	//channel.Close(err) // 关闭连接

	//clientProxy := proxy.NewRemoteClientProxy(service.MyService())
	//clientProxy.Invoke("Hello", nil)

}

func TestServer(t *testing.T) {
	nacosServer := server.NewNacosServerStarter("127.0.0.1:9527", "127.0.0.1:8848")

	err := nacosServer.Start()

	if err != nil {
		log.Fatal("start nacosServer failed:", err)
	}

}

func TestClient1(t *testing.T) {
	// 子连接的流水线配置
	var clientnitializer = func(channel netty.Channel) {
		channel.Pipeline().
			// 最大允许包长128字节，使用\n分割包, 丢弃分隔符
			AddLast(frame.DelimiterCodec(1024, "\r\n", true)).
			//AddLast(frame.LengthFieldCodec(binary.BigEndian, 2048, 8, 2, 12, 0)).
			AddLast(netty.ReadIdleHandler(time.Second*3), netty.WriteIdleHandler(time.Second*5)).
			AddLast(codec.CommonCodec(0, 8, 71, 1)).
			AddLast(client.NewClientHandler())
	}

	channel, err1 := netty.NewBootstrap(netty.WithClientInitializer(clientnitializer)).Connect(":9528")
	if err1 != nil {
		log.Fatal("channel err: ", err1)
		return
	}

	parameters := []interface{}{"hello，this is go language"}

	message := protocol.RpcRequestProtocol("123455", "helloService", "sayHello", parameters,
		[]string{"java.lang.String"}, "java.lang.String", false, "1.0.1", false)

	err := channel.Write(message)
	if err != nil {
		log.Fatal("channel1 err: ", err)
	}

	time.Sleep(time.Second * 120)

}

func TestServer1(t *testing.T) {

	// 子连接的流水线配置
	var childInitializer = func(channel netty.Channel) {
		channel.Pipeline().
			// 最大允许包长128字节，使用\n分割包, 丢弃分隔符
			AddLast(frame.DelimiterCodec(1024, "\r\n", true)).
			AddLast(netty.ReadIdleHandler(time.Second * 30)).
			//AddLast(frame.LengthFieldCodec(binary.BigEndian, 2048, 8, 2, 12, 0)).
			AddLast(codec.CommonCodec(0, 8, 71, 1)).
			AddLast(server.NewServerHandler())
	}

	// ErrServerClosed is returned by the Server call Shutdown or Close.
	//var ErrServerClosed = errors.New("netty: Server closed")

	// 创建Bootstrap & 监听端口 & 接受连接
	bs := netty.NewBootstrap(netty.WithChildInitializer(childInitializer))

	//bs.Listen(":9527", tcp.WithOptions(tcpOptions)).Async(func(err error) {
	//	if nil != err && ErrServerClosed != err {
	//		t.Fatal(err)
	//	}
	//})
	err := bs.Listen(":9527").Sync()
	if err != nil {
		log.Fatal("LISTEN ERR: ", err)
	}
	bs.Shutdown()
}
