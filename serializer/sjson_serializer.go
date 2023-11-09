package serializer

/**
Hessian 序列化 实现 公共序列化 接口
*/

func NewSJsonSerializer() CommonSerializer {
	return &sjsonSerializer{

		Value: 6,
	}
}

/*
*

	跨 语言 序列化协议
*/
type sjsonSerializer struct {
	Value int
}

func (sjsonSerializer *sjsonSerializer) Serialize(message any) (data []byte, err error) {
	data, err = json.Marshal(message)
	return
}

func (sjsonSerializer *sjsonSerializer) Deserialize(data []byte, message interface{}) (target any, err error) {
	err = json.Unmarshal(data, message)
	target = message
	return
}

func (sjsonSerializer *sjsonSerializer) GetValue() int {
	return sjsonSerializer.Value
}
