package serializer

const (
	KryoSerializerCode    = 0
	JsonSerializerCode    = 1
	HessianSerializerCode = 2
	FurySerializerCode    = 3
	// 处理请求包 - Java（客户端） -> Golang（服务端） 跨协议 JSON 序列化、Golang（服务端） -> Java（客户端） 跨协议 JSON 反序列化
	CJsonSerializerCode = 6
	// 处理响应包 - Golang（客户端） -> Java（服务端） 跨协议 JSON 反序列化、Java（服务端） -> Golang（客户端） 跨协议 JSON 序列化
	SJsonSerializerCode = 7
)

/*
*

	公共序列化接口
*/
type CommonSerializer interface {
	Serialize(message any) (data []byte, err error)

	Deserialize(data []byte, message interface{}) (target any, err error)

	GetValue() int
}
