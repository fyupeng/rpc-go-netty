package serializer

/**
Kryo 序列化 实现 公共序列化 接口
*/

func NewKryoSerializer() CommonSerializer {
	return &kryoSerializer{
		Value: 0,
	}
}

type kryoSerializer struct {
	Value int
}

func (kryoSerializer *kryoSerializer) Serialize(message any) (data []byte, err error) {
	//TODO
	return
}

func (kryoSerializer *kryoSerializer) Deserialize(data []byte, message interface{}) (target interface{}, err error) {
	//TODO
	return
}

func (kryoSerializer *kryoSerializer) GetValue() int {
	return kryoSerializer.Value
}
