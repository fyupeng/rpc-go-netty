package serializer

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

/**
Json 序列化 实现 公共序列化 接口
*/

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func NewJsonSerializer() CommonSerializer {
	return &jsonSerializer{
		Value: 1,
	}
}

type jsonSerializer struct {
	Value int
}

func (jsonSerializer *jsonSerializer) Serialize(message any) (data []byte, err error) {
	data, err = json.Marshal(message)
	return
}

func (jsonSerializer *jsonSerializer) Deserialize(data []byte, message interface{}) (target any, err error) {
	fmt.Println(message)
	err = json.Unmarshal(data, message)
	target = message
	fmt.Println("message ", message)
	return
}

func (jsonSerializer *jsonSerializer) GetValue() int {
	return jsonSerializer.Value
}
