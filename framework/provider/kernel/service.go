package kernel

import (
	"hade/framework/gin"
	"net/http"
)

// 引擎服务
type HadeKernelService struct {
	engine *gin.Engine
}

// NewHadeKernelService 创建hade
func NewHadeKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &HadeKernelService{engine: httpEngine}, nil
}

// 返回web引擎
func (s *HadeKernelService) HttpEngine() http.Handler {
	return s.engine
}
