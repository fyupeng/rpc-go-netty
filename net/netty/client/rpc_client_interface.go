package client

import (
	"rpc-go-netty/net/netty/future"
	"rpc-go-netty/protocol"
)

type RpcClient interface {
	SendRequest(rpcRequest protocol.RequestProtocol) *future.CompleteFuture
}
