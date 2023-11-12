package client

type RpcClient interface {
	SendRequest(interfaceName string, methodName string, parameters []interface{}, paramTypes []string, returnTypes []string)
}
