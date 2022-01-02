package config

import (
	"hade/framework"
	"hade/framework/contact"
	"hade/framework/provider/app"
	"path/filepath"
)

type HadeConfigProvider struct {
	framework.ServiceProvider

	Folder string
}

// Register registe a new function for make a service instance
func (provider *HadeConfigProvider) Register(c framework.Container) framework.NewInstance {
	return NewHadeConfig
}

// Boot will called when the service instantiate
func (provider *HadeConfigProvider) Boot(c framework.Container) error {
	c.InstanceMustExist(contact.AppKey)
	c.InstanceMustExist(contact.EnvKey)

	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *HadeConfigProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *HadeConfigProvider) Params(c framework.Container) []interface{} {
	appServicer := c.MustGetInstance(contact.AppKey).(app.HadeApp)
	envServicer := c.MustGetInstance(contact.EnvKey).(contact.Env)
	env := envServicer.AppEnv()
	// 配置文件夹地址
	configFolder := appServicer.ConfigFolder()
	envFolder := filepath.Join(configFolder, env)
	return []interface{}{c, envFolder, envServicer.All()}
}

/// Name define the name for this service
func (provider *HadeConfigProvider) Name() string {
	return contact.ConfigKey
}
