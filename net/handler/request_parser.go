package handler

import (
	"fmt"
	"github.com/go-netty/go-netty"
	"log"
	"reflect"
	"rpc-go-netty/protocol"
	"rpc-go-netty/serializer"
)

func NewRequestParser(serializerCode int) netty.ChannelHandler {
	return &requestParser{
		serializerCode: serializerCode,
	}
}

type requestParser struct {
	idleEvent      int32
	serializerCode int // 序列化类型
}

func (h *requestParser) HandleActive(ctx netty.ActiveContext) {
	//TODO implement me
	ctx.HandleActive()
}

func (req *requestParser) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	//TODO implement me
	log.Println("prepare for golang protocol transfer to java:: ", message)

	request := message.(protocol.RequestProtocol)
	request.SetParamTypes(adaptParamTypes(req.serializerCode, request.GetParamTypes()))
	request.SetReturnType(adaptReturnTypes(req.serializerCode, request.GetReturnType()))
	request.SetMethodName(adaptJavaMethodName(req.serializerCode, request.GetMethodName()))

	ctx.HandleWrite(message)
}

func (req *requestParser) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	//TODO implement me
	ctx.HandleInactive(ex)
}

func (req *requestParser) HandleRead(ctx netty.InboundContext, message netty.Message) {

	log.Println("protocolHandler receive message from server: ", message)

	ctx.HandleRead(message)
}

func (req *requestParser) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	// 处理异常情况
	log.Println("Exception caught:", ex)
	//ctx.HandleException(ex)
}

func adaptParamTypes(serializerCode int, paramTypes []string) []string {
	adaptParamTypes := make([]string, len(paramTypes))
	if serializer.CJsonSerializerCode == serializerCode {
		for index, param := range paramTypes {
			switch param {
			case reflect.String.String():
				adaptParamTypes[index] = "java.lang.String"
			case reflect.Float64.String():
				adaptParamTypes[index] = "double"
			case reflect.Float32.String():
				adaptParamTypes[index] = "float"
			case reflect.Int.String():
				adaptParamTypes[index] = "int"
			}
		}
		return adaptParamTypes
	} else {
		return paramTypes
	}

}

func adaptReturnTypes(serializerCode int, returnType string) string {
	var adaptReturnType string
	if serializer.CJsonSerializerCode == serializerCode {
		switch returnType {
		case reflect.String.String():
			adaptReturnType = "java.lang.String"
		case reflect.Float64.String():
			adaptReturnType = "double"
		case reflect.Float32.String():
			adaptReturnType = "float"
		case reflect.Int.String():
			adaptReturnType = "int"
		}
		return adaptReturnType
	} else {
		return adaptReturnType
	}

}

func adaptJavaMethodName(serializerCode int, methodName string) string {
	// Java 请求远程方法调用适配
	if serializer.CJsonSerializerCode == serializerCode {
		return Capitalize(methodName)
	}
	return methodName
}

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
