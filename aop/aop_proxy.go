package aop

import (
	"log"
	"reflect"
	"rpc-go-netty/discovery/service_discovery"
	"rpc-go-netty/protocol"
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

func NewClientProxy(serviceConsumer service_discovery.ServiceDiscovery) Proxy {
	return &clientProxy{
		ServiceConsumer: serviceConsumer,
	}
}

type clientProxy struct {
	ServiceConsumer service_discovery.ServiceDiscovery
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

	serviceAddr, getServiceErr := proxy.ServiceConsumer.LookupServiceWithGroupName(interfaceName, "1.0.1")

	if getServiceErr != nil {
		log.Fatal("get Service Fatal: ", getServiceErr)
	}

	channel := proxy.ServiceConsumer.GetChannel(serviceAddr.String())

	message := protocol.RpcRequestProtocol("123455", interfaceName, methodName, parameters,
		paramTypes, returnTypes[0], false, "1.0.1", false)

	err := channel.Write(message)
	if err != nil {
		log.Fatal("channel1 err: ", err)
	}

	time.Sleep(time.Second * 100)

}

func getServiceName(pkgPath string, reveiverName string) string {
	pkgs := strings.Split(pkgPath, "/")
	var serviceName string
	for i := 1; i < len(pkgs); i++ {
		serviceName += pkgs[i] + "."
	}
	return serviceName + reveiverName
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
