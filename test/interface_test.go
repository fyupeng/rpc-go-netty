package test

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInterface(t *testing.T) {
	typeByName := reflect.TypeOf((*interface{})(nil)).Elem()
	interfaceType, found := typeByName.FieldByName("cn/fyupeng/cn.fyupeng.service.HelloWorldServiceImpl")
	if !found {
		fmt.Errorf("Interface type not found")
	}
	fmt.Println(interfaceType)
}
