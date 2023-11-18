package handler

import (
	"github.com/go-netty/go-netty"
	"log"
	"reflect"
	"rpc-go-netty/factory"
	"rpc-go-netty/net/netty/future"
	"rpc-go-netty/protocol"
)

func NewClientHandler() netty.ChannelHandler {
	return &clientHandler{
		unProcessResult: (factory.GetInstance(reflect.TypeOf((*future.UnProcessResult)(nil)))).(*future.UnProcessResult),
	}
}

type clientHandler struct {
	idleEvent       int32
	unProcessResult *future.UnProcessResult
}

func (h *clientHandler) HandleActive(ctx netty.ActiveContext) {
	//TODO implement me
	ctx.HandleActive()
}

func (h *clientHandler) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	//TODO implement me
	log.Println("client send request to server: ", message)
	ctx.HandleWrite(message)
}

func (h *clientHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	//TODO implement me
	ctx.HandleInactive(ex)
}

func (h *clientHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {

	log.Println("client receive message from server: ", message)

	response := message.(protocol.ResponseProtocol)

	h.unProcessResult.Complete(response.GetRequestId(), response)

	//ctx.HandleRead(message)
}

func (h *clientHandler) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	// 处理异常情况
	log.Println("Exception caught:", ex)
	//ctx.HandleException(ex)
}

func (h *clientHandler) HandleEvent(ctx netty.EventContext, event netty.Event) {

	switch event.(type) {
	case netty.ReadIdleEvent:
		break
	case netty.WriteIdleEvent:
		log.Println("Send heartbeat packets to server[%s]", ctx.Channel().RemoteAddr())
		ctx.Write(protocol.RpcRequestProtocol("", "", "", []interface{}{}, []string{}, "", false, "", true))
	default:
		panic(event)
	}
}
