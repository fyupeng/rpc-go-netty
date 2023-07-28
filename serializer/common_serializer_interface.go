package serializer

/*
*

	公共序列化接口
*/
type CommonSerializer interface {
	Serialize(message any) (data []byte, err error)

	Deserialize(data []byte, message any) (err error)

	GetValue() int
}
