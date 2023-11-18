package cn_fyupeng_service

type HelloWorldService interface {
	SayHello(message string) string
}

type HelloWorldServiceImpl struct {
	name  string `annotation:"cn.fyupeng.service.HelloWorldService"`
	group string `annotation:"1.0.0"`
}

func (helloService *HelloWorldServiceImpl) SayHello(message string) string {
	return "GoLang say: I have receive your message: " + message
}
