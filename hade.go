package hade

import (
	"hade/framework"
	"hade/framework/gin"
	"net/http"
)

type Hade struct {
	// Container 容器
	Container framework.Container
	// Engine web引擎
	Engine *gin.Engine
	// server http服务
	Server http.Server
}

// New 初始化Hade实例
func New() Hade {
	hade := Hade{
		Container: framework.NewHadeContainer(),
		Engine:    gin.New(),
	}

	return hade
}

// Bind 为Hade容器绑定一个服务
func (hade *Hade) Bind(service framework.ServiceProvider) {
	hade.boot()
	hade.Container.Bind(service)
	hade.Engine.SetContainer(hade.Container)
}

// Use 为Hade添加中间件
func (hade *Hade) Use(middleware ...gin.HandlerFunc) {
	hade.boot()
	for _, v := range middleware {
		hade.Engine.Use(v)
	}
}

// Run 启动Hade服务
func (hade *Hade) Run(addr string) {
	server := &http.Server{
		Handler: hade.Engine,
		Addr:    "127.0.0.1:9200",
	}
	hade.Server = *server

	server.ListenAndServe()
}

// boot 是否初始化一个Hade实例
func (hade *Hade) boot() {
	if hade.Container == nil {
		panic("hade container is not init")
	}
	if hade.Engine == nil {
		panic("hade engin is not init")
	}
}
