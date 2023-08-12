package protocol

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
	return rpcRequestProtocol.RequestId
}

func (rpcRequestProtocol *rpcRequestProtocol) GetInterfaceName() string {
	//TODO implement me
	return rpcRequestProtocol.InterfaceName
}

func (rpcRequestProtocol *rpcRequestProtocol) GetMethodName() string {
	return rpcRequestProtocol.MethodName
}

func (rpcRequestProtocol *rpcRequestProtocol) GetParameters() []interface{} {
	return rpcRequestProtocol.Parameters
}

func (rpcRequestProtocol *rpcRequestProtocol) GetHeartBeat() bool {
	return rpcRequestProtocol.HeartBeat
}

type rpcRequestProtocol struct {
	RequestId     string        `json:"requestId" hessian:"requestId"`
	InterfaceName string        `json:"interfaceName" hessian:"interfaceName"`
	MethodName    string        `json:"methodName" hessian:"methodName"`
	Parameters    []interface{} `json:"parameters" hessian:"parameters"`
	ParamTypes    []string      `json:"paramTypes" hessian:"paramTypes"`
	ReturnType    string        `json:"returnType" hessian:"returnType"`
	ReSend        bool          `json:"reSend" hessian:"reSend"`
	Group         string        `json:"group" hessian:"group"`
	HeartBeat     bool          `json:"heartBeat" hessian:"heartBeat"`
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

func (proto rpcResponseProtocol) Ok(requestId, message string) Protocol {
	proto.RequestId = requestId
	proto.Message = message
	proto.StatusCode = 200
	return proto
}

func (proto rpcResponseProtocol) Success(requestId string, data interface{}) Protocol {
	proto.RequestId = requestId
	proto.Message = "ok"
	proto.StatusCode = 200
	proto.Data = data
	return proto
}

func (proto rpcResponseProtocol) SuccessWithCheckCode(requestId string, data interface{}, checkCode []byte) Protocol {
	proto.RequestId = requestId
	proto.Message = "ok"
	proto.StatusCode = 200
	proto.Data = data
	proto.CheckCode = checkCode
	return proto
}

type rpcResponseProtocol struct {
	RequestId  string      `json:"requestId"`
	CheckCode  []byte      `json:"checkCode"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
