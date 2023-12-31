package test

import (
	"fmt"
	"github.com/fyupeng/rpc-go-netty/annotation"
	"github.com/fyupeng/rpc-go-netty/aop"
	cn_fyupeng_service "github.com/fyupeng/rpc-go-netty/cn.fyupeng.service"
	"github.com/fyupeng/rpc-go-netty/discovery/load_balancer"
	"github.com/fyupeng/rpc-go-netty/discovery/service_discovery"
	"github.com/fyupeng/rpc-go-netty/net/netty/client"
	"github.com/fyupeng/rpc-go-netty/net/netty/server"
	"github.com/fyupeng/rpc-go-netty/protocol"
	"github.com/fyupeng/rpc-go-netty/serializer"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestProxy(t *testing.T) {
	// 直连 方式（不使用 注册中心 和 负载均衡）
	//client := client.NewNettyClientDirect("192.168.67.191:9527", serializer.JsonSerializerCode)
	// 通过 注册中心负载获取
	client := client.NewNettyClient2Alone(load_balancer.NewRandLoadBalancer(), serializer.JsonSerializerCode, "127.0.0.1:8848")
	h := aop.NewClientProxy(client)
	h.AddAnnotation(&annotation.Annotation{
		GroupName:   "1.0.0",
		ServiceName: "cn.fyupeng.service.HelloWorldService",
	})
	h.Invoke(reflect.TypeOf((*cn_fyupeng_service.HelloWorldService)(nil)), "SayHello", []interface{}{"这是go代理端"})
}
func TestClient(t *testing.T) {

	serviceConsumer := service_discovery.NewServiceConsumerAlone(load_balancer.NewRandLoadBalancer(), serializer.JsonSerializerCode, "127.0.0.1:8848")

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
	// 创建 服务端 和 连接 注册中心
	nacosServer := server.NewNettyServer("192.168.67.191:9527", "127.0.0.1:8848", serializer.JsonSerializerCode)
	// 发布 服务
	nacosServer.PublishService(&cn_fyupeng_service.HelloWorldServiceImpl{})
	// 启动服务 并 监听 客户端请求
	err := nacosServer.Start()
	if err != nil {
		log.Fatal("start nacosServer failed:", err)
	}
}
