package server

type ServerStart interface {
	Start() (err error)

	publishService()
}
