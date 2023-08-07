package serializer

import (
	hessian "github.com/apache/dubbo-go-hessian2"
)

/**
Hessian 序列化 实现 公共序列化 接口
*/

func NewHessianSerializer() CommonSerializer {
	return &hessianSerializer{
		Value: 2,
	}
}

type hessianSerializer struct {
	Value int
}

func (hessianSerializer *hessianSerializer) Serialize(message any) (data []byte, err error) {
	// 创建一个编码器
	encoder := hessian.NewEncoder()
	// 将数据编码为字节数组
	err = encoder.Encode(message)
	data = encoder.Buffer()
	encoder.Clean()
	return
}

func (hessianSerializer *hessianSerializer) Deserialize(data []byte, message interface{}) (target any, err error) {
	decoder := hessian.NewDecoder(data)

	target, err = decoder.DecodeValue()
	decoder.Clean()
	return
}

func (hessianSerializer *hessianSerializer) GetValue() int {
	return hessianSerializer.Value
}
