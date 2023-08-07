package service

type HelloWorldService interface {
	sayHello(message string) string
}

type HelloWorldServiceImpl struct {
}

func (helloService *HelloWorldServiceImpl) sayHello(message string) string {
	return "I say: " + message
}
