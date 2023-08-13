package service_provider

type ServiceProvider interface {
	GetService(serviceName string) (service interface{})
	AddService(service interface{}, serviceName string)
}
