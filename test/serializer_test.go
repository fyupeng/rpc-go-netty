package test

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log"
	"rpc-go-netty/serializer"
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
