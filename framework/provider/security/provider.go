package security

import (
	"hade/framework"
	"hade/framework/contact"
	"time"

	"github.com/spf13/cast"
)

type HadeSecurityProvider struct {
	framework.ServiceProvider

	duration time.Duration
	slat     string
}

// Register 注册方法
func (provider *HadeSecurityProvider) Register(container framework.Container) framework.NewInstance {
	return NewSecurityServicer
}

// Boot 启动调用
func (provider *HadeSecurityProvider) Boot(container framework.Container) error {
	configServicer := container.MustGetInstance(contact.ConfigKey).(contact.Config)

	if provider.duration == 0 && configServicer.IsExist("security.duration") {
		provider.duration = time.Duration(cast.ToInt64(configServicer.GetString("security.duration")))
	} else if provider.duration == 0 {
		panic("duration is not exist")
	}

	if provider.slat == "" && configServicer.IsExist("security.slat") {
		provider.slat = configServicer.GetString("security.slat")
	} else if provider.slat == "" {
		panic("slat is not exist")
	}

	return nil
}

// IsDefer 是否延迟初始化
func (provider *HadeSecurityProvider) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (provider *HadeSecurityProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, provider.duration, provider.slat}
}

// Name 获取字符串凭证
func (provider *HadeSecurityProvider) Name() string {
	return contact.SecurityKey
}
