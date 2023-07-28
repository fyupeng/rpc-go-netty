package protocol

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
	GetHeartBeat() bool
}

/*
*
响应协议
*/
type ResponseProtocol interface {
	Protocol
}
