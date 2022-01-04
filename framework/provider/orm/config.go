package orm

import (
	"context"
	"hade/framework"
	"hade/framework/contact"

	"gorm.io/gorm"
)

// GetBaseConfig 读取database.yaml根目录结构
func GetBaseConfig(c framework.Container) map[string]*contact.DBConfig {
	configService := c.MustGetInstance(contact.ConfigKey).(contact.Config)
	logService := c.MustGetInstance(contact.LogKey).(contact.Log)

	configMap := make(map[string]*contact.DBConfig, 0)
	databaseMap := configService.GetStringMap("database")
	for k, v := range databaseMap {
		config := &contact.DBConfig{}
		// 直接使用配置服务的load方法读取,yaml文件
		err := configService.Load("database."+k, config, v)
		if err != nil {
			// 直接使用logService来打印错误信息
			logService.Error(context.Background(), "parse database config error", nil)
			return nil
		}
		configMap[k] = config
	}

	return configMap
}

// WithConfigPath 加载配置文件地址
func WithConfigPath(configPath string) contact.DBOption {
	return func(container framework.Container, config *contact.DBConfig) error {
		configService := container.MustGetInstance(contact.ConfigKey).(contact.Config)
		// 加载configPath配置路径
		if err := configService.Load(configPath, config, nil); err != nil {
			return err
		}
		return nil
	}
}

// WithGormConfig 表示自行配置Gorm的配置信息
func WithGormConfig(gormConfig *gorm.Config) contact.DBOption {
	return func(container framework.Container, config *contact.DBConfig) error {
		if gormConfig.Logger == nil {
			gormConfig.Logger = config.Logger
		}
		config.Config = gormConfig
		return nil
	}
}

// WithDryRun 设置空跑模式
func WithDryRun() contact.DBOption {
	return func(container framework.Container, config *contact.DBConfig) error {
		config.DryRun = true
		return nil
	}
}

// WithFullSaveAssociations 设置保存时候关联
func WithFullSaveAssociations() contact.DBOption {
	return func(container framework.Container, config *contact.DBConfig) error {
		config.FullSaveAssociations = true
		return nil
	}
}
