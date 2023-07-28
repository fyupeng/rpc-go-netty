package proxy

import (
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
	"reflect"
)

// MyService 是一个远程服务接口
type MyService interface {
	Hello(name string) string
}

// MyServiceImpl 是 MyService 的具体实现
type MyServiceImpl struct{}

// Hello 是 MyService 接口定义的方法的具体实现
func (s *MyServiceImpl) Hello(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func NewRemoteClientProxy(serviceInterface interface{}) ClientProxy {
	serviceInterfaceType := GetInterfaceType(serviceInterface)
	clientProxy := &RemoteClientProxy{
		serviceImplPackage: "net/proxy/service",
	}
	// 一个接口只能有一个实现，多个实现只返回第一个实现
	serviceImpl := scanPackage(clientProxy.serviceImplPackage, serviceInterfaceType)
	proxyHandler := reflect.ValueOf(clientProxy)
	serviceImpl.String()
	for i := 0; i < serviceInterfaceType.NumMethod(); i++ {
		method := serviceInterfaceType.Method(i)
		name := method.Name
		// 接口所有方法 都由 Invoke 方法代理调用
		proxyHandler.MethodByName("Invoke").Call([]reflect.Value{
			reflect.ValueOf(name),
			reflect.ValueOf([]reflect.Value{}),
			// 代理调用实际方法处理
			reflect.MakeFunc(method.Type, func(args []reflect.Value) (results []reflect.Value) {
				// 处理 序列化、编解码
				return reflect.ValueOf(serviceImpl).Method(i).Call(args)
			}),
		})

	}
	return &RemoteClientProxy{}
}

type RemoteClientProxy struct {
	serviceImplPackage string
}

func (clientProxy *RemoteClientProxy) Invoke(methodName string, args []interface{}) interface{} {
	return nil
}

func scanPackage(packageName string, serviceInterface reflect.Type) reflect.Value {

	implPackage := packageName
	fmt.Println("implPackage" + implPackage)
	defs, err := getPkgDefinitions(implPackage)
	if err != nil {
		fmt.Println("Failed to get package definitions:", err)
	}

	for _, def := range defs {
		if reflect.PtrTo(def).Implements(serviceInterface) {
			val := reflect.New(def).Elem()
			return val
		}
	}
	return reflect.Value{}
}

func getPkgDefinitions(pkg string) ([]reflect.Type, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.LoadTypes,
	}, pkg)
	if err != nil {
		return nil, err
	}

	fmt.Println("pkgs ")
	fmt.Println(pkgs)

	var defs []reflect.Type
	for _, pkg := range pkgs {
		fmt.Println(pkg)
		fmt.Println(pkg.Name)
		fmt.Println(pkg.PkgPath)
		fmt.Println(pkg.TypesInfo)
		for _, obj := range pkg.TypesInfo.Defs {
			fmt.Println(obj)
			if obj, ok := obj.(*types.TypeName); ok {
				typ := obj.Type().Underlying()
				defs = append(defs, reflect.TypeOf(typ))
			}
		}
	}

	return defs, nil
}

func GetInterfaceType(i interface{}) reflect.Type {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
