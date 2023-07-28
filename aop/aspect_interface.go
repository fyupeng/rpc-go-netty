package aop

type AspectInterface interface {
	Before(point *JoinPoint) bool
	After(point *JoinPoint)
	Finally(point *JoinPoint)
	GetAspectExpress() string
}
