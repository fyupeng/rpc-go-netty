package protocol

import (
	"fmt"
)

/*
*
RPC 请求协议
*/
func RpcRequestProtocol(requestId string, interfaceName string, methodName string, parameters []interface{},
	paramTypes []string, returnType string, reSend bool, group string, heartBeat bool) RequestProtocol {
	return &rpcRequestProtocol{
		RequestId:     requestId,
		InterfaceName: interfaceName,
		MethodName:    methodName,
		Parameters:    parameters,
		ParamTypes:    paramTypes,
		ReturnType:    returnType,
		ReSend:        reSend,
		Group:         group,
		HeartBeat:     heartBeat,
	}
}

func NewRpcRequestProtocol() RequestProtocol {
	return &rpcRequestProtocol{}
}

func (rpcRequestProtocol *rpcRequestProtocol) GetRequestId() string {
	fmt.Println("111111")
	return rpcRequestProtocol.RequestId
}

func (rpcRequestProtocol *rpcRequestProtocol) GetHeartBeat() bool {
	return rpcRequestProtocol.HeartBeat
}

type rpcRequestProtocol struct {
	RequestId     string        `json:"requestId"`
	InterfaceName string        `json:"interfaceName"`
	MethodName    string        `json:"methodName"`
	Parameters    []interface{} `json:"parameters"`
	ParamTypes    []string      `json:"paramTypes"`
	ReturnType    string        `json:"returnType"`
	ReSend        bool          `json:"reSend"`
	Group         string        `json:"group"`
	HeartBeat     bool          `json:"heartBeat"`
}

/*
*
RPC 响应协议
*/
func RpcResponseProtocol(requestId string, checkCode []byte, statusCode int, message string, data interface{}) ResponseProtocol {
	return &rpcResponseProtocol{
		RequestId:  requestId,
		CheckCode:  checkCode,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

func NewRpcResponseProtocol() ResponseProtocol {
	return &rpcResponseProtocol{}
}

type rpcResponseProtocol struct {
	RequestId  string      `json:"requestId"`
	CheckCode  []byte      `json:"checkCode"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
