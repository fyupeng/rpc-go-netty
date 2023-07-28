package serializer

/**
Hessian 序列化 实现 公共序列化 接口
*/

func HessianSerializer() CommonSerializer {
	return &jsonSerializer{
		Value: 2,
	}
}

type hessianSerializer struct {
	Value int
}

func (hessianSerializer *hessianSerializer) Serialize(message any) (data []byte, err error) {
	return
}

func (hessianSerializer *hessianSerializer) Deserialize(data []byte, message any) (err error) {
	return
}

func (hessianSerializer *hessianSerializer) GetValue() int {
	return hessianSerializer.Value
}
