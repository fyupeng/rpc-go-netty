package serializer

import (
	"fmt"
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

	log.Println("reqMap")
	log.Println(reqMap)

	if err = mapstructure.Decode(reqMap, message); err != nil {

		fmt.Println(err)

	}
	fmt.Println(message)
	target = message

	//res := "{"
	//for k, v := range mapJson[0] {
	//	res += "\"" + k + "\"" + ":" + v + ","
	//}
	//// 除去多余 ,
	//res = res[:len(res)-1]
	//res += "}"
	//
	//log.Println("res")
	//log.Println(res)
	//
	//err = json.Unmarshal([]byte(res), message)
	//
	//target = message

	return
}

func (cjsonSerializer *cjsonSerializer) GetValue() int {
	return cjsonSerializer.Value
}
