package handler

import (
	"github.com/go-netty/go-netty"
	"log"
	"rpc-go-netty/protocol"
)

func NewResponseParser(serializerCode int) netty.ChannelHandler {
	return &responseParser{
		serializerCode: serializerCode,
	}
}

type responseParser struct {
	idleEvent      int32
	serializerCode int // 序列化类型
}

func (resp *responseParser) HandleActive(ctx netty.ActiveContext) {
	//TODO implement me
	ctx.HandleActive()
}

func (resp *responseParser) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	//TODO implement me
	log.Println("prepare for golang protocol transfer to java:: ", message)
	response := message.(protocol.ResponseProtocol)
	ctx.HandleWrite(response)
}

func (h *responseParser) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	//TODO implement me
	ctx.HandleInactive(ex)
}

func (h *responseParser) HandleRead(ctx netty.InboundContext, message netty.Message) {

	log.Println("protocolHandler receive message from server: ", message)

	ctx.HandleRead(message)
}

func (h *responseParser) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	// 处理异常情况
	log.Println("Exception caught:", ex)
	//ctx.HandleException(ex)
}
