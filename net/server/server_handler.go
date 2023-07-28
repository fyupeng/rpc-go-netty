package server

import (
	"errors"
	"fmt"
	"github.com/go-netty/go-netty"
	"log"
	"rpc-go-netty/protocol"
	"sync/atomic"
)

func NewServerHandler() netty.ChannelHandler {
	return &serverHandler{}
}

type serverHandler struct {
	idleEvent int32
}

func (h *serverHandler) HandleActive(ctx netty.ActiveContext) {
	//TODO implement me
	ctx.HandleActive()
}

func (h *serverHandler) HandleWrite(ctx netty.OutboundContext, message netty.Message) {
	//TODO implement me
	fmt.Println("HandleWrite")
	fmt.Println("服务端触发写操作：这是写操作信息： ", message)
	//ctx.Write(message)
	ctx.HandleWrite(message)
}

func (h *serverHandler) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	//TODO implement me
	ctx.HandleInactive(ex)
}

func (h *serverHandler) HandleRead(ctx netty.InboundContext, message netty.Message) {
	// 在这里处理已经解码的消息
	// 根据您的具体业务逻辑进行处理
	// 示例：打印消息内容
	requestProtocol := message.(protocol.RequestProtocol)
	if requestProtocol.GetHeartBeat() {
		log.Println("server receive heartBeat packets")
		return
	}
	fmt.Println("服务端读到客户端消息： ", message)

	// 示例：发送响应消息
	//response := "Hello, Client!"
	//parameters := []interface{}{"hello，这里是go语言"}
	//message = protocol.RpcRequestProtocol("123455", "helloService", "sayHello", parameters,
	//	[]string{"java.lang.String"}, "java.lang.String", false, "1.0.0", false)
	//ctx.Write(message)
	// 最后一个处理器不用将消息交给下一个
	//ctx.HandleRead(message)
}

func (h *serverHandler) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	// 处理异常情况
	fmt.Println("Exception caught:", ex)
	ctx.Close(ex)
	//ctx.HandleException(ex)
}

func (h *serverHandler) HandleEvent(ctx netty.EventContext, event netty.Event) {

	switch event.(type) {
	case netty.ReadIdleEvent:
		fmt.Println("read idle event", " isActive: ", ctx.Channel().IsActive())
		if 2 == atomic.AddInt32(&h.idleEvent, 1) {
			ctx.Channel().Close(errors.New("headrbeat packets hve not been received for a long time, server close the channel with client"))
		}

	default:
		panic(event)
	}
}
