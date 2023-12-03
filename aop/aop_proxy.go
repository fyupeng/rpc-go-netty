package aop

import (
	"fmt"
	"github.com/fyupeng/rpc-go-netty/annotation"
	"github.com/fyupeng/rpc-go-netty/net/netty/client"
	"github.com/fyupeng/rpc-go-netty/protocol"
	"github.com/fyupeng/rpc-go-netty/utils/idworker"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func init() {
	RegisterPoint(reflect.TypeOf((*clientProxy)(nil)))
	RegisterAspect(&Aspect{})
}

type Aspect struct {
}

func (a *Aspect) Before(point *JoinPoint) bool {
	//fmt.Println("before")
	return true
}

func (a *Aspect) After(point *JoinPoint) {
	//fmt.Println("after")
}

func (a *Aspect) Finally(point *JoinPoint) {
	//fmt.Println("finally")
}

func (a *Aspect) GetAspectExpress() string {
	return ".*\\.Invoke"
}

type Proxy interface {
	Invoke(interfaceType reflect.Type, methodName string, parameters []interface{})
	AddAnnotation(annotation *annotation.Annotation)
}

func NewClientProxy(nettyClient client.RpcClient) Proxy {
	worker, err := idworker.NewNowWorker(0)
	if err != nil {
		log.Fatal("id create failed: ", err)
	}
	return &clientProxy{
		nettyClient: nettyClient,
		IdWorker:    worker,
	}
}

type clientProxy struct {
	nettyClient client.RpcClient
	IdWorker    idworker.IdWorker
	annot       *annotation.Annotation
}

func (proxy *clientProxy) AddAnnotation(annotation *annotation.Annotation) {
	proxy.annot = annotation
}

func (proxy *clientProxy) Invoke(interfaceType reflect.Type, methodName string, parameters []interface{}) {

	pkgPath := interfaceType.PkgPath()
	receiverName := interfaceType.Name()
	if interfaceType.Kind() == reflect.Ptr {
		pkgPath = interfaceType.Elem().PkgPath()
		receiverName = interfaceType.Elem().Name()
	}

	interfaceName := getServiceName(pkgPath, receiverName)

	method := getInterfaceMethod(interfaceType, methodName)

	methodType := method.Type

	if len(parameters) != methodType.NumIn() {
		panic("parameter length not equal the paramters for method[" + method.Name + "]")
	}

	paramTypes := make([]string, methodType.NumIn())
	for i := 0; i < len(paramTypes); i++ {
		paramTypes[i] = methodType.In(i).String()
	}

	returnTypes := make([]string, methodType.NumOut())
	for i := 0; i < len(returnTypes); i++ {
		returnTypes[i] = methodType.Out(i).String()
	}

	// 封装成sendRequest
	requestId := strconv.FormatInt(proxy.IdWorker.NextId(), 10)

	groupName := ""

	if proxy.annot.GroupName != "" || proxy.annot.ServiceName != "" {
		if proxy.annot.GroupName != "" {
			groupName = proxy.annot.GroupName
		}
		if proxy.annot.ServiceName != "" {
			interfaceName = proxy.annot.ServiceName
		}
	}

	rpcRequest := protocol.RpcRequestProtocol(requestId, interfaceName, methodName, parameters, paramTypes, returnTypes[0], false, groupName, false)

	completeFuture := proxy.nettyClient.SendRequest(rpcRequest)

	futureResult, err := completeFuture.GetFuture()

	if err != nil {
		log.Fatal("ClientProxy Invoke failed:", err)
	}

	fmt.Println(futureResult)

	time.Sleep(time.Second * 10)

}

func getServiceName(pkgPath string, receiverName string) string {
	pkgs := strings.Split(pkgPath, "/")
	var serviceName string
	for i := 1; i < len(pkgs); i++ {
		serviceName += pkgs[i] + "."
	}
	return serviceName + receiverName
}

func getInterfaceMethod(interfaceType reflect.Type, methodName string) reflect.Method {
	method, isExist := interfaceType.MethodByName(methodName)
	if interfaceType.Kind() == reflect.Ptr {
		method, isExist = interfaceType.Elem().MethodByName(methodName)
	}
	if isExist {
		return method
	}
	panic(methodName + "this method [" + methodName + "] is not exist in interface " + interfaceType.Name())
}
