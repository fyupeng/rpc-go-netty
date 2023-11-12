package test

import (
	"fmt"
	"log"
	"reflect"
	"rpc-go-netty/aop"
	cn_fyupeng_service "rpc-go-netty/cn.fyupeng.service"
	"rpc-go-netty/discovery/load_balancer"
	"rpc-go-netty/discovery/service_discovery"
	"rpc-go-netty/net/netty/client"
	"rpc-go-netty/net/netty/server"
	"rpc-go-netty/protocol"
	"rpc-go-netty/serializer"
	"testing"
	"time"
)

func TestProxy(t *testing.T) {
	client := client.NewNettyClient2Alone(load_balancer.NewRandLoadBalancer(), serializer.CJsonSerializerCode, "127.0.0.1:8848")
	h := aop.NewClientProxy(client)
	h.Invoke(reflect.TypeOf((*cn_fyupeng_service.HelloWorldService)(nil)), "SayHello", []interface{}{"这是go代理端"})
}
func TestClient(t *testing.T) {

	serviceConsumer := service_discovery.NewServiceConsumer(load_balancer.NewRandLoadBalancer(), serializer.JsonSerializerCode, "127.0.0.1:8848")

	//serviceAddr, getServiceErr := serviceConsumer.LookupService("TestService")
	//
	//if getServiceErr != nil {
	//	log.Fatal("get Service Fatal: ", getServiceErr)
	//}

	serviceAddr, getServiceErr := serviceConsumer.LookupServiceWithGroupName("helloService", "1.0.1")

	fmt.Println(serviceAddr)

	if getServiceErr != nil {
		log.Fatal("get Service Fatal: ", getServiceErr)
	}

	channel := serviceConsumer.GetChannel("127.0.0.1:9527")

	parameters := []interface{}{"hello，这里是go语言"}

	message := protocol.RpcRequestProtocol("123455", "helloService", "sayHello", parameters,
		[]string{"java.lang.String"}, "java.lang.String", false, "1.0.1", false)

	err := channel.Write(message)
	if err != nil {
		log.Fatal("channel1 err: ", err)
	}

	time.Sleep(time.Second * 60)

	// 在这里可以进行其他操作，发送消息等

	//channel.Close(err) // 关闭连接

	//clientProxy := proxy.NewRemoteClientProxy(cn.fyupeng.service.MyService())
	//clientProxy.Invoke("Hello", nil)

}

func TestServer(t *testing.T) {
	nacosServer := server.NewNettyServer("192.168.43.33:9527", "127.0.0.1:8848", serializer.SJsonSerializerCode)

	nacosServer.PublishService(&cn_fyupeng_service.HelloWorldServiceImpl{})

	err := nacosServer.Start()

	if err != nil {
		log.Fatal("start nacosServer failed:", err)
	}

}
