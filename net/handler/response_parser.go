package handler

import (
	"fmt"
	"github.com/go-netty/go-netty"
	"log"
	"reflect"
	"rpc-go-netty/protocol"
	"rpc-go-netty/serializer"
	"rpc-go-netty/utils/aes"
)

func NewResponseParser(serializerCode int) netty.ChannelHandler {
	return &responseParser{
		serializerCode: serializerCode,
		jsonSerializer: serializer.NewFJsonSerializer(),
	}
}

type responseParser struct {
	idleEvent      int32
	serializerCode int // 序列化类型
	jsonSerializer serializer.CommonSerializer
}

func (resp *responseParser) HandleActive(ctx netty.ActiveContext) {
	//TODO implement me
	ctx.HandleActive()
}

func (resp *responseParser) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	//TODO implement me
	response := message.(protocol.ResponseProtocol)

	// GoLang 获取的 返回值以多个数组形式返回，需要支持 Java
	// 后面重试机制 优化 拷贝 方式
	result := response.GetData()
	response.SetData(result)
	response.SetDataType(adaptDataType(resp.serializerCode, response.GetDataType()))

	checkData, err := resp.jsonSerializer.Serialize(result)
	if err != nil {
		log.Fatal(fmt.Sprintf("get checkData serializer failed,origin: %v, after: %s", result, checkData))
	}

	checkCode, err := aes.Encrypt(string(checkData))
	if err != nil {
		log.Fatal(fmt.Sprintf("get checkCode Encrypt failed,origin: %v, after: %s", checkData, checkCode))
	}

	response.SetCheckCode(checkCode)

	// goLang 的 base64 是 以 "_" 作为分隔符，而 java 以 / 作为分隔符
	if serializer.SJsonSerializerCode == resp.serializerCode {
		log.Println("prepare response for golang protocol transfer to java:", message)
	} else {
		log.Println("prepare response for golang protocol transfer to golang:", message)
	}
	ctx.HandleWrite(message)
}

func (h *responseParser) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	//TODO implement me
	ctx.HandleInactive(ex)
}

func (h *responseParser) HandleRead(ctx netty.InboundContext, message netty.Message) {

	ctx.HandleRead(message)
}

func (h *responseParser) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	// 处理异常情况
	log.Println("Exception caught:", ex)
	//ctx.HandleException(ex)
}

func adaptDataType(serializerCode int, dataType string) string {
	var adaptDataType string
	if serializer.SJsonSerializerCode == serializerCode {
		switch dataType {
		case reflect.String.String():
			adaptDataType = "java.lang.String"
		case reflect.Float64.String():
			adaptDataType = "double"
		case reflect.Float32.String():
			adaptDataType = "float"
		case reflect.Int.String():
			adaptDataType = "int"
		}
		return adaptDataType
	} else {
		return dataType
	}

}
