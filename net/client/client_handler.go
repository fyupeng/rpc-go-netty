package client

import (
	"fmt"
	"github.com/go-netty/go-netty"
	"log"
	"rpc-go-netty/protocol"
)

func NewClientHandler() netty.ChannelHandler {
	return &clientHandler{}
}

type clientHandler struct {
	idleEvent int32
}

func (h *clientHandler) HandleActive(ctx netty.ActiveContext) {
	//TODO implement me
	fmt.Println("->", "active:", ctx.Channel().RemoteAddr())

	// 给对端发送一条消息，将进入如下流程（视编解码配置）
	// Text -> TextCodec -> LengthFieldCodec   -> Channel.Write
	// 文本     文本编码      组装协议格式（长度字段）     网络发送
	ctx.HandleActive()
}

func (h *clientHandler) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	//TODO implement me
	fmt.Println("HandleWrite")
	fmt.Println("客户端触发写操作：这是写操作信息： ", message)
	ctx.HandleWrite(message)
}

func (h *clientHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	//TODO implement me
	fmt.Println("->", "inactive:", ctx.Channel().RemoteAddr(), ex)
	ctx.HandleInactive(ex)
}

func (h *clientHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {

	fmt.Println("clientHandler: ", message)

	//ctx.HandleRead(message)
}

func (h *clientHandler) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	// 处理异常情况
	fmt.Println("Exception caught:", ex)
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
