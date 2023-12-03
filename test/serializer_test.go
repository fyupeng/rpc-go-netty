package test

import (
	"encoding/json"
	"fmt"
	"github.com/fyupeng/rpc-go-netty/protocol"
	"github.com/fyupeng/rpc-go-netty/serializer"
	jsoniter "github.com/json-iterator/go"
	"log"
	"testing"
)

func TestSerializer(t *testing.T) {
	hessianSerializer := serializer.NewHessianSerializer()
	jsonSerializer := serializer.NewJsonSerializer()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	student := &Student{
		Name: "小明",
		Age:  18,
		Sex:  true,
	}
	fmt.Println("student: ", student)
	data, err := hessianSerializer.Serialize(student)
	data1, err1 := jsonSerializer.Serialize(student)
	data3, err3 := json.Marshal(student)

	if err != nil {
		log.Fatal("err: ", err)
	}
	fmt.Println("serializered data: ", data)

	if err1 != nil {
		log.Fatal("err: ", err1)
	}
	fmt.Println("serializered data1: ", data1)

	if err3 != nil {
		log.Fatal("err3: ", err3)
	}
	fmt.Println("serializered data3: ", data3)

	var decodData any
	decodData, err = hessianSerializer.Deserialize(data, &Student{})

	var decodData1 any
	decodData1, err1 = jsonSerializer.Deserialize(data1, &Student{})

	fmt.Println("decodData: ", decodData)
	fmt.Println("decodData1: ", decodData1)

	obj, ok := decodData.(*Student)
	if ok {
		fmt.Println("decodData is ok", obj)
	}
	obj1, ok1 := decodData1.(*Student)
	if ok1 {
		fmt.Println("decodData is ok", obj1)
	}

}

type people interface {
}

func NewPeople(idx int) people {
	switch idx {
	case 0:
		return &student{}
	case 1:
		return &teacher{}
	}
	return &student{}
}

type student struct {
	Name    string
	Age     int
	Sex     bool
	Tuition int
}

type teacher struct {
	Name   string
	Age    int
	Sex    bool
	Salary int
}

func TestStudentDeserialized(t *testing.T) {
	student := student{
		Name:    "小明",
		Age:     18,
		Sex:     true,
		Tuition: 8500,
	}
	studentBinary, err := json.Marshal(student)
	if err != nil {
		log.Fatal("序列化失败！", err)
	}
	fmt.Println("序列化成功：", studentBinary)
	// 动态获取请求
	idx := 0
	message := NewPeople(idx)
	// 动态选择 people 序列化
	err = json.Unmarshal(studentBinary, message)
	if err != nil {
		log.Fatal("反序列化失败", err)
	}
	fmt.Println("反序列化成功：", message)
}

func getProtocolByCode(protocolCode int) (proto protocol.Protocol) {
	switch protocolCode {
	case protocol.RequestProtocolCode:
		proto = protocol.NewRpcRequestProtocol()
	case protocol.ResponseProtocolCode:
		proto = protocol.NewRpcResponseProtocol()
	case protocol.UnRecognizeProtocolCode:
		log.Println("unrecognized protocol:", protocolCode)
	default:
		log.Println("unrecognized protocol:", protocolCode)
	}
	return
}

func TestProtocolDeserialized(t *testing.T) {
	//proto := protocol.RpcRequestProtocol("123455", "HelloService", "sayHello", []interface{}{"arg0", "arg1"},
	//	[]string{"string", "string"}, "string", false, "1.0.1", false)
	proto := protocol.RpcResponseProtocol("123455", "VYGThW0MXPf4v88IKP/o4g==", 200, "测试消息", "这是服务端消息")
	protocolBinary, err := json.Marshal(proto)
	if err != nil {
		log.Fatal("序列化失败！", err)
	}
	fmt.Println("序列化成功：", protocolBinary)
	// 动态获取请求
	idx := protocol.ResponseProtocolCode
	message := getProtocolByCode(idx)
	// 动态选择 people 序列化
	err = json.Unmarshal(protocolBinary, message)
	if err != nil {
		log.Fatal("反序列化失败", err)
	}
	fmt.Println("反序列化成功：", message)
}

func TestNullDataDeserialzed(t *testing.T) {

}
