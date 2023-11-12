package aop

import (
	"reflect"
	"rpc-go-netty/net/netty/client"
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
}

func NewClientProxy(nettyClient client.RpcClient) Proxy {
	return &clientProxy{
		nettyClient: nettyClient,
	}
}

type clientProxy struct {
	nettyClient client.RpcClient
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

	proxy.nettyClient.SendRequest(interfaceName, methodName, parameters, paramTypes, returnTypes)

	time.Sleep(time.Second * 100)

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
