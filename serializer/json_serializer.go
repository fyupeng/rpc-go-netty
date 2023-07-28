package serializer

import "encoding/json"

/**
Json 序列化 实现 公共序列化 接口
*/

func JsonSerializer() CommonSerializer {
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

func (jsonSerializer *jsonSerializer) Deserialize(data []byte, message any) (err error) {
	err = json.Unmarshal(data, message)
	return
}

func (jsonSerializer *jsonSerializer) GetValue() int {
	return jsonSerializer.Value
}
