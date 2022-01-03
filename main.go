package hade

// func main() {
// 	type Hade struct {
// 		// Container 容器
// 		Container framework.Container
// 		// Handler 路由handler
// 		Handler http.Handler
// 	}

// 	fmt.Println("===")
// 	fmt.Println(Hade{})

// 	var baseFolder string
// 	flag.StringVar(&baseFolder, "base_folder", "", "base_folder参数, 默认为当前路径")
// 	flag.Parse()

// 	container := framework.NewHadeContainer()
// 	// 绑定App服务提供者
// 	container.Bind(&app.HadeAppProvider{BaseFolder: baseFolder})
// 	// 后续初始化需要绑定的服务提供者...
// 	container.Bind(&env.HadeEnvProvider{})
// 	container.Bind(&config.HadeConfigProvider{})
// 	container.Bind(&hadeLog.HadeLogServicerProvider{})

// 	engine, err := core.RunHttpEngine(container)
// 	if err != nil {
// 		container.Bind(&kernel.HadeKernelProvider{HttpEngine: engine})
// 	}

// 	server := &http.Server{
// 		Handler: engine,
// 		Addr:    "127.0.0.1:9200",
// 	}

// 	go func() {
// 		server.ListenAndServe()
// 	}()

// 	quit := make(chan os.Signal, 1)

// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

// 	<-quit

// 	if err := server.Shutdown(context.Background()); err != nil {
// 		log.Fatal("Server ShutDown", err.Error())
// 	}
// }
