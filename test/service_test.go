package test

import (
	"fmt"
	"reflect"
	"testing"
)

func TestService(t *testing.T) {
	Register(&Student{}, &Student{})
}

func TestSlice(t *testing.T) {
	args1 := "123123"
	args11 := reflect.ValueOf(args1).Slice(0, len(args1))
	fmt.Println(args11)

	args2 := []interface{}{"12312", "12312312"}
	args22 := reflect.ValueOf(args2).Slice(0, len(args2))
	fmt.Println(args22)

	args3 := "123123"
	args33 := reflect.ValueOf(args3)
	fmt.Println(args33)
}

func Register(services ...interface{}) {
	for _, service := range services {
		fmt.Println(service)

		serviceType := reflect.TypeOf(service)
		fmt.Println(serviceType)

		fmt.Println(serviceType.String())
		fmt.Println(serviceType.String())

		fmt.Println(serviceType.Elem())
		fmt.Println(serviceType.Elem().Name())

		// Check if the method "SayHello" exists
		sayHelloMethod, isValid := serviceType.MethodByName("SayHello")
		if isValid {
			arg0 := "hello"

			// Create a new instance of the cn.fyupeng.service type
			serviceValue := reflect.New(serviceType.Elem())
			interfaceValue := serviceValue.Interface()

			result := sayHelloMethod.Func.Call([]reflect.Value{reflect.ValueOf(interfaceValue), reflect.ValueOf(arg0)})
			fmt.Println(result)
		} else {
			fmt.Println("SayHello method not found")
		}
	}
}

// 定义接口
type Shape interface {
	Area() float64
}

// 定义结构体
type Rectangle struct {
	Length float64
	Width  float64
}

// 让结构体实现接口中定义的方法
func (r Rectangle) Area() float64 {
	return r.Length * r.Width
}

func TestInterfaceService(t *testing.T) {
	// 创建结构体的实例
	rect := Rectangle{Length: 5, Width: 3}

	// 获取结构体的反射类型
	rectType := reflect.TypeOf(rect)

	// 遍历结构体实现的接口列表

	rectType.String()

	for i := 0; i < rectType.NumField(); i++ {
		field := rectType.Field(i)
		if field.Type.Kind() == reflect.Interface {
			fmt.Printf("Implements interface: %s\n", field.Type.Name())
		}
	}
}
