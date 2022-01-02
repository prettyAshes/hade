package kernel

import (
	"errors"
	"hade/framework"
	"hade/framework/contact"
	"hade/framework/gin"
)

// HadeKernelProvider 提供web引擎
type HadeKernelProvider struct {
	HttpEngine *gin.Engine
}

func (provider *HadeKernelProvider) Name() string {
	return contact.KernelKey
}

func (provider *HadeKernelProvider) Register(container framework.Container) framework.NewInstance {
	return NewHadeKernelService
}

func (provider *HadeKernelProvider) IsDefer() bool {
	return false
}

func (provider *HadeKernelProvider) Params(container framework.Container) []interface{} {
	return []interface{}{provider.HttpEngine}
}

func (provider *HadeKernelProvider) Boot(container framework.Container) error {
	if provider.HttpEngine == nil {
		return errors.New("HttpEngine is nil")
	}
	return nil
}
