package client

import "rpc-go-netty/net/netty/future"

type RpcClient interface {
	SendRequest(interfaceName string, methodName string, parameters []interface{}, paramTypes []string, returnTypes []string) *future.CompleteFuture
}
