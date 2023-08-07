package serializer

const (
	KryoSerializerCode    = 0
	JsonSerializerCode    = 1
	HessianSerializerCode = 2
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
