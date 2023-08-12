package cn_fyupeng_service

type HelloWorldService interface {
	sayHello(message string) string
	Haha(message string) string
}

type HelloWorldServiceImpl struct {
	name  string `annotation:"cn.fyupeng.service.HelloWorldService"`
	group string `annotation:"1.0.0"`
}

func (helloService *HelloWorldServiceImpl) sayHello(message string) string {
	return "I say: " + message
}

func (helloService *HelloWorldServiceImpl) Haha(message string) string {
	return "I say: " + message
}
