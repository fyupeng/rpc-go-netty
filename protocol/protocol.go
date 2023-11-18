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

func (rpcRequestProtocol *rpcRequestProtocol) SetMethodName(methodName string) {
	rpcRequestProtocol.MethodName = methodName
}

func (rpcRequestProtocol *rpcRequestProtocol) GetParameters() []interface{} {
	return rpcRequestProtocol.Parameters
}

func (rpcRequestProtocol *rpcRequestProtocol) GetParamTypes() []string {
	return rpcRequestProtocol.ParamTypes
}

func (rpcRequestProtocol *rpcRequestProtocol) SetParamTypes(paramTypes []string) {
	//TODO implement me
	rpcRequestProtocol.ParamTypes = paramTypes
}

func (rpcRequestProtocol *rpcRequestProtocol) GetHeartBeat() bool {
	return rpcRequestProtocol.HeartBeat
}

func (rpcRequestProtocol *rpcRequestProtocol) GetReturnType() string {
	return rpcRequestProtocol.ReturnType
}

func (rpcRequestProtocol *rpcRequestProtocol) SetReturnType(returnType string) {
	//TODO implement me
	rpcRequestProtocol.ReturnType = returnType
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
func RpcResponseProtocol(requestId string, checkCode string, statusCode int, message string, data interface{}) *rpcResponseProtocol {
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

func (rpcResponseProtocol *rpcResponseProtocol) Ok(requestId, message string) Protocol {
	rpcResponseProtocol.RequestId = requestId
	rpcResponseProtocol.Message = message
	rpcResponseProtocol.StatusCode = 200
	return rpcResponseProtocol
}

func (rpcResponseProtocol *rpcResponseProtocol) Success(requestId string, data interface{}) Protocol {
	rpcResponseProtocol.RequestId = requestId
	rpcResponseProtocol.Message = "ok"
	rpcResponseProtocol.StatusCode = 200
	rpcResponseProtocol.Data = data
	return rpcResponseProtocol
}

func (rpcResponseProtocol *rpcResponseProtocol) SuccessWithCheckCode(requestId string, data interface{}, dataType string, checkCode string) Protocol {
	rpcResponseProtocol.RequestId = requestId
	rpcResponseProtocol.Message = "ok"
	rpcResponseProtocol.StatusCode = 200
	rpcResponseProtocol.Data = data
	rpcResponseProtocol.DataType = dataType
	rpcResponseProtocol.CheckCode = checkCode
	return rpcResponseProtocol
}

func (rpcResponseProtocol *rpcResponseProtocol) GetRequestId() string {
	return rpcResponseProtocol.RequestId

}

func (rpcResponseProtocol *rpcResponseProtocol) GetData() interface{} {
	return rpcResponseProtocol.Data
}

func (rpcResponseProtocol *rpcResponseProtocol) SetData(data interface{}) {
	rpcResponseProtocol.Data = data
}

func (rpcResponseProtocol *rpcResponseProtocol) GetCheckCode() string {
	return rpcResponseProtocol.CheckCode
}

func (rpcResponseProtocol *rpcResponseProtocol) SetCheckCode(checkCode string) {
	rpcResponseProtocol.CheckCode = checkCode
}

func (rpcResponseProtocol *rpcResponseProtocol) GetDataType() string {
	return rpcResponseProtocol.DataType
}

func (rpcResponseProtocol *rpcResponseProtocol) SetDataType(dataType string) {
	rpcResponseProtocol.DataType = dataType
}

type rpcResponseProtocol struct {
	RequestId  string      `json:"requestId" hessian:"requestId"`
	CheckCode  string      `json:"checkCode" hessian:"checkCode"`
	StatusCode int         `json:"statusCode" hessian:"statusCode"`
	Message    string      `json:"message" hessian:"message"`
	Data       interface{} `json:"data" hessian:"data"`
	DataType   string      `json:"dataType" hessian:"dataType"`
}
