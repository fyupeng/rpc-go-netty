package test

type Student struct {
	Name string
	Age  int
	Sex  bool
}

func (Student) JavaClassName() string {
	return "test.Student"
}
