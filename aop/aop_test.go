package aop

import (
	"fmt"
	"reflect"
	"testing"
)

func init() {
	RegisterPoint(reflect.TypeOf((*HelloAop)(nil)))
	RegisterAspect(&Aspect{})
}

type Aspect struct{}

func (a *Aspect) Before(point *JoinPoint) bool {
	fmt.Println("before")
	return true
}

func (a *Aspect) After(point *JoinPoint) {
	fmt.Println("after")
}

func (a *Aspect) Finally(point *JoinPoint) {
	fmt.Println("finally")
}

func (a *Aspect) GetAspectExpress() string {
	return ".*\\.HelloAop"
}

type HelloAop struct {
}

func (h *HelloAop) HelloAop() {
	fmt.Println("helloAop")
}

func TestAop(t *testing.T) {
	h := &HelloAop{}
	h.HelloAop()
}
