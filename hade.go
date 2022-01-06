package hade

import (
	"hade/framework"
	"hade/framework/contact"
	"hade/framework/gin"
	"hade/middlewares"
	"net/http"

	hadeLog "hade/framework/provider/log"
	"hade/framework/provider/orm"

	"hade/framework/provider/app"
	"hade/framework/provider/config"
	"hade/framework/provider/env"
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
func New(baseFolder string) Hade {
	hade := Hade{
		Container: framework.NewHadeContainer(),
		Engine:    gin.New(),
	}
	hade.bind(baseFolder)

	logService := hade.Container.MustGetInstance(contact.LogKey).(contact.Log)
	hade.Use(gin.Recovery())
	hade.Use(gin.LoggerWithWriter(logService.GetOutPut()))
	hade.Use(middlewares.Logger(hade.Container))

	return hade
}

// bind 绑定常用服务
func (hade *Hade) bind(baseFolder string) {
	// 绑定App服务提供者
	hade.Bind(&app.HadeAppProvider{BaseFolder: baseFolder})
	// 后续初始化需要绑定的服务提供者...
	hade.Bind(&env.HadeEnvProvider{})
	hade.Bind(&config.HadeConfigProvider{})
	hade.Bind(&hadeLog.HadeLogServicerProvider{})
	hade.Bind(&orm.GormProvider{})
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
