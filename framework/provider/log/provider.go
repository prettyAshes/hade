package log

import (
	"hade/framework"
	"hade/framework/contact"
	"hade/framework/provider/log/formatter"
	"hade/framework/provider/log/services"
	"io"

	"github.com/spf13/cast"
)

// HadeLogServicerProvider 提供日志服务
type HadeLogServicerProvider struct {
	framework.ServiceProvider

	Driver string

	// 日志级别
	Level contact.LogLevel
	// 日志输出格式方法
	Formatter contact.Formatter
	// 日志context上下文信息获取函数
	CtxFielder contact.CtxFielder
	// 日志输出信息
	Output io.Writer
	// 日志路径
	LogFolder string
	// 日期格式
	DateFormat string
}

func (provider *HadeLogServicerProvider) Name() string {
	return contact.LogKey
}

func (provider *HadeLogServicerProvider) Boot(container framework.Container) error {
	configServicer := container.MustGetInstance(contact.ConfigKey).(contact.Config)

	if configServicer.IsExist("log.driver") {
		provider.Driver = configServicer.GetString("log.driver")
	} else {
		panic("log driver is not exist")
	}

	if configServicer.IsExist("log.formatter") {
		v := configServicer.GetString("log.formatter")
		if v == "json" {
			provider.Formatter = formatter.JsonFormatter
		} else if v == "text" {
			provider.Formatter = formatter.TextFormatter
		}
	} else {
		provider.Formatter = formatter.TextFormatter
	}

	logLevel := cast.ToUint32(configServicer.Get("log.level"))
	if logLevel == uint32(contact.UnknownLevel) {
		provider.Level = contact.InfoLevel
	}

	return nil
}

func (provider *HadeLogServicerProvider) Register(container framework.Container) framework.NewInstance {
	// 根据driver的配置项确定
	switch provider.Driver {
	case "single":
		return services.NewHadeSingleLog
	case "rotate":
		return services.NewHadeRotateLog
	case "console":
		return services.NewHadeConsoleLog
	case "custom":
		return services.NewHadeCustomLog
	default:
		return services.NewHadeConsoleLog
	}
}

func (provider *HadeLogServicerProvider) IsDefer() bool {
	return false
}

func (provider *HadeLogServicerProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, provider.Level, provider.CtxFielder, provider.Formatter, provider.Output}
}
