package conn

type Router interface {
	// 在处理conn业务之前的钩子方法
	PreHandle(request *Request)
	// 在处理conn业务的方法
	Handle(request *Request)
	// 处理conn业务之后的钩子方法
	PostHandle(request *Request)
}

// 基类先继承完接口方法，子类可以只实现部分方法
type BaseRouter struct{}

func (b *BaseRouter) PreHandle(request *Request)  {}
func (b *BaseRouter) Handle(request *Request)     {}
func (b *BaseRouter) PostHandle(request *Request) {}

var _ Router = (*BaseRouter)(nil)
