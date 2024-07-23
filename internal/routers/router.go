package routers

// IRouter ...
type IRouter[T any] interface {
	Setup(app T)
}
