package serializer

import "fmt"

/**
Hessian 序列化 实现 公共序列化 接口
*/

func NewFurySerializer() CommonSerializer {
	return &furySerializer{

		Value: 3,
	}
}

type furySerializer struct {
	Value int
}

func (furySerializer *furySerializer) Serialize(message any) (data []byte, err error) {
	// 创建一个编码器
	_ = fmt.Errorf("unsuported serializer!")
	return
}

func (furySerializer *furySerializer) Deserialize(data []byte, message interface{}) (target any, err error) {
	_ = fmt.Errorf("unsuported serializer!")
	return
}

func (furySerializer *furySerializer) GetValue() int {
	return furySerializer.Value
}
