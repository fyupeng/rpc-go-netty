package protocol

const (
	RequestProtocolCode     = 71
	ResponseProtocolCode    = 73
	UnRecognizeProtocolCode = -1
)

/*
*
协议
*/
type Protocol interface {
}

/*
*
请求协议
*/
type RequestProtocol interface {
	Protocol
	GetRequestId() string
	GetInterfaceName() string
	GetMethodName() string
	SetMethodName(methodName string)
	GetParameters() []interface{}
	GetParamTypes() []string
	SetParamTypes(paramTypes []string)
	GetHeartBeat() bool
	GetReturnType() string
	SetReturnType(returnType string)
}

/*
*
响应协议
*/
type ResponseProtocol interface {
	Protocol
	Ok(requestId, message string) Protocol
	Success(requestId string, data interface{}) Protocol
	SuccessWithCheckCode(requestId string, data interface{}, dataType string, checkCode string) Protocol

	GetData() interface{}
	SetData(data interface{})
	GetCheckCode() string
	SetCheckCode(checkCode string)
	GetDataType() string
	SetDataType(dataType string)
}
