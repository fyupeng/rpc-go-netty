package demo

import "fmt"

// @Register("service0", "arg1", "arg2")
type MyStruct1 struct {
	Name string
}

func (s *MyStruct1) Print() {
	fmt.Println("This is", s.Name)
}

// @Register("service1", "arg1", "arg2")
func MyFunction1(arg1 string, arg2 string) (string, error) {
	fmt.Println("This is a function.")
	return "", nil
}

// @Register("service2", "arg1", "arg2")
func MyFunction2(arg1 string, arg2 string) error {
	fmt.Println("This is a function.")
	return nil
}
