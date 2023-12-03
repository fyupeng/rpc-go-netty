package client

import (
	"github.com/fyupeng/rpc-go-netty/net/netty/future"
	"github.com/fyupeng/rpc-go-netty/protocol"
)

type RpcClient interface {
	SendRequest(rpcRequest protocol.RequestProtocol) *future.CompleteFuture
}
