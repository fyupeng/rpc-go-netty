package server

type RpcServer interface {
	Start() (err error)

	PublishService(services ...interface{})
}
