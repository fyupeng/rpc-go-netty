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
	var result []reflect.Value
	if len(args) == 1 {
		result = method.Func.Call([]reflect.Value{interfaceValue, reflect.ValueOf(args[0])})
	} else {
		// 调用方法
		result = method.Func.Call([]reflect.Value{interfaceValue, reflect.ValueOf(args).Slice(0, len(args))})
	}

	// 处理返回值
	resultVal := reflect.ValueOf(result)

	if resultVal.Kind() != reflect.Slice {
		// go/java 返回首个 具体数据类型数据
		return result[0].Interface()
	} else if resultVal.Len() > 0 {
		return result[0].Interface()
	} else {
		return result
	}

	return nil
}
