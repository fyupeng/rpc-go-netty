package serializer

import (
	"github.com/goinggo/mapstructure"
	"log"
)

/**
Hessian 序列化 实现 公共序列化 接口
*/

func NewFJsonSerializer() CommonSerializer {
	return &cjsonSerializer{

		Value: 6,
	}
}

/*
*

	跨 语言 序列化协议
*/
type cjsonSerializer struct {
	Value int
}

func (cjsonSerializer *cjsonSerializer) Serialize(message any) (data []byte, err error) {
	data, err = json.Marshal(message)
	return
}

func (cjsonSerializer *cjsonSerializer) Deserialize(data []byte, message interface{}) (target any, err error) {
	var reqMap map[string]interface{}
	err = json.Unmarshal(data, &reqMap)

	if err = mapstructure.Decode(reqMap, message); err != nil {
		log.Println(err)
	}
	target = message

	return
}

func (cjsonSerializer *cjsonSerializer) GetValue() int {
	return cjsonSerializer.Value
}
