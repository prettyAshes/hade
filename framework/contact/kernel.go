package contact

import "net/http"

const KernelKey = "hade:kernel"

// Kernel 项目启动服务接口
type Kernel interface {
	// HttpEngin 获取http.Handler
	HttpEngin() http.Handler
}
