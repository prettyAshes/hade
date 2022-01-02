package framework

type NewInstance func(...interface{}) (interface{}, error)

// ServiceProvider 定义一个服务提供者需要实现的功能
type ServiceProvider interface {
	// Name 服务名称
	Name() string
	// Register 服务注册到容器中
	Register(container Container) NewInstance
	// IsDefer 是否延迟
	IsDefer() bool
	// Params 服务参数
	Params(container Container) []interface{}
	// Boot 容器是否做好准备工作
	Boot(container Container) error
}
