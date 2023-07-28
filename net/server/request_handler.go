package server

import "reflect"

func NewRequestHandler() InvocationHandler {
	return &requestHandler{}
}

type requestHandler struct {
}

func (handle *requestHandler) Handle(methodName string, args []interface{}) interface{} {
	// 构造远程方法
	method := reflect.ValueOf(handle).MethodByName(methodName)
	if !method.IsValid() {
		return nil
	}

	// 构造参数列表
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	// 调用方法
	result := method.Call(in)

	// 提取结果
	if len(result) > 0 {
		return result[0].Interface()
	}

	return nil
}
