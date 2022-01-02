package framework

import (
	"errors"
	"sync"
)

// Container 定义一个容器需要实现的功能
type Container interface {
	// Bind 容器绑定一个服务
	Bind(service ServiceProvider) error
	// IsBind 判断服务是否绑定在容器中
	IsBind(name string) bool
	// Make 获取一个服务
	Make(name string) (interface{}, error)
	// MakeNew 重新初始化一个服务
	MakeNew(name string) (interface{}, error)
	// InstanceMustExist 服务实例必须存在，用于provider中Boot判断
	InstanceMustExist(name string)
	// MustGetInstance 获取一个服务实例, 优先从内存中获取
	MustGetInstance(name string) interface{}
}

// HadeContainer hade框架结构体
type HadeContainer struct {
	Container
	// Providers 服务
	providers map[string]ServiceProvider
	// Instances 服务实例
	instances map[string]interface{}
	// Lock 读写锁
	lock sync.RWMutex
}

// NewHadeContainer 创建一个服务容器
func NewHadeContainer() *HadeContainer {
	return &HadeContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

// Bind 容器绑定一个服务
func (hade *HadeContainer) Bind(serviceProvider ServiceProvider) error {
	serviceName := serviceProvider.Name()

	hade.providers[serviceName] = serviceProvider

	if !serviceProvider.IsDefer() {
		instance, err := hade.newInstance(serviceProvider)
		if err != nil {
			return err
		}
		hade.instances[serviceName] = instance
	}

	return nil
}

// IsBind 服务是否绑定
func (hade *HadeContainer) IsBind(name string) bool {
	return hade.findServiceProvider(name) != nil
}

// Make 获取一个服务
func (hade *HadeContainer) Make(name string) (interface{}, error) {
	return hade.make(name, false)
}

// Make 获取一个服务
func (hade *HadeContainer) MakeNew(name string) (interface{}, error) {
	return hade.make(name, true)
}

// InstanceMustExist 服务实例必须存在
func (hade *HadeContainer) InstanceMustExist(name string) {
	if _, ok := hade.instances[name]; !ok {
		panic(name + " instance is not exist")
	}
}

// MustGetInstance 获取一个服务实例, 优先从内存中获取
func (hade *HadeContainer) MustGetInstance(name string) interface{} {
	hade.lock.RLock()
	defer hade.lock.RUnlock()

	if instance, ok := hade.instances[name]; ok {
		return instance
	}

	serviceProvider := hade.findServiceProvider(name)
	if serviceProvider == nil {
		panic(name + " servicer is not find")
	}

	instance, err := hade.make(name, true)
	if err != nil {
		panic(err)
	}

	return instance
}

// newInstance 创建一个服务实例
func (hade *HadeContainer) newInstance(serviceProvider ServiceProvider) (interface{}, error) {
	if err := serviceProvider.Boot(hade); err != nil {
		return nil, err
	}
	serviceParams := serviceProvider.Params(hade)
	method := serviceProvider.Register(hade)
	instance, err := method(serviceParams...)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

// findServiceProvider 查找一个服务提供者
func (hade *HadeContainer) findServiceProvider(serviceName string) ServiceProvider {
	hade.lock.RLock()
	defer hade.lock.RUnlock()

	if servicerProvider, ok := hade.providers[serviceName]; ok {
		return servicerProvider
	}

	return nil
}

// make 获取一个服务
func (hade *HadeContainer) make(name string, forceNew bool) (interface{}, error) {
	hade.lock.Lock()

	if hade.providers[name] == nil {
		return nil, errors.New("contract " + name + " have not register")
	}

	serviceProvider := hade.providers[name]

	hade.lock.Unlock()

	if forceNew {
		return hade.newInstance(serviceProvider)
	}

	if hade.instances[name] != nil {
		instance, err := hade.newInstance(serviceProvider)
		if err != nil {
			return nil, err
		}

		hade.instances[name] = instance
	}

	return hade.instances[name], nil
}
