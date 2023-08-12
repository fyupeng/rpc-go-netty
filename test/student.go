package test

import "fmt"

type Student struct {
	Name string
	Age  int
	Sex  bool
}

func (stu Student) JavaClassName() string {
	return "cn.fyupeng.service.Student"
}

func (stu Student) AutoRegister() {

}

func (stu Student) SayHello(message string) {
	fmt.Println("say hello: ", message)
}
