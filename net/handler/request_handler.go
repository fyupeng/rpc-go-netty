package handler

import (
	"log"
	"reflect"
)

func NewRequestHandler() InvocationHandler {
	return &requestHandler{}
}

type requestHandler struct {
}

func (handle *requestHandler) Handle(service interface{}, methodName string, args []interface{}) interface{} {
	// 构造远程方法

	serviceType := reflect.TypeOf(service)

	method, isValid := serviceType.MethodByName(methodName)

	if !isValid {
		log.Fatal("methodName valid:", isValid)
	}

	// 创建接口实例
	interfaceValue := reflect.New(serviceType).Elem()

	// 调用函数
	var result interface{}
	if len(args) == 1 {
		result = method.Func.Call([]reflect.Value{interfaceValue, reflect.ValueOf(args[0])})
	} else {
		// 调用方法
		result = method.Func.Call([]reflect.Value{interfaceValue, reflect.ValueOf(args).Slice(0, len(args))})
	}

	// 处理返回值
	resultVal := reflect.ValueOf(result)

	if resultVal.Kind() != reflect.Slice {
		// go/java 直接返回
		return result
	} else if resultVal.Len() > 0 {
		// java 模式，有多个返回值默认返回一个

		// go 模式 直接返回
		return result
	}

	return nil
}
