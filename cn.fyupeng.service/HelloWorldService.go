package cn_fyupeng_service

type HelloWorldService interface {
	Haha(message string) string
}

type HelloWorldServiceImpl struct {
	name  string `annotation:"cn.fyupeng.service.HelloWorldService"`
	group string `annotation:"1.0.0"`
}

func (helloService *HelloWorldServiceImpl) Haha(message string) string {
	return "I say: " + message
}
