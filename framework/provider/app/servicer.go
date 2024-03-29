package app

import (
	"errors"
	"hade/framework"
	"hade/framework/util"
	"path/filepath"

	"github.com/google/uuid"
)

// HadeApp 代表hade框架的App实现
type HadeApp struct {
	container  framework.Container // 服务容器
	baseFolder string              // 基础路径
	appId      string              // 表示当前这个app的唯一id, 可以用于分布式锁等
	configMap  map[string]string   // 配置加载
	config     AppConfig           // 配置
}

// AppConfig App配置文件的配置
type AppConfig struct {
	Swagger bool `yaml:"swagger"`
	Path    struct {
		ConfigFolder     string `yaml:"config_folder"`
		LogFolder        string `yaml:"log_folder"`
		HttpFolder       string `yaml:"http_folder"`
		ConsoleFolder    string `yaml:"console_folder"`
		StorageFolder    string `yaml:"storage_folder"`
		ProviderFolder   string `yaml:"provider_folder"`
		MiddlewareFolder string `yaml:"middleware_folder"`
		CommandFolder    string `yaml:"command_folder"`
		RuntimeFolder    string `yaml:"runtime_folder"`
		TestFolder       string `yaml:"test_folder"`
		DeployFolder     string `yaml:"deploy_folder"`
		AppFolder        string `yaml:"app_folder"`
	} `yaml:"path"`
	Dev struct {
		Port int `yaml:"port"`
	} `yaml:"dev"`
}

// AppID 表示这个App的唯一ID
func (app HadeApp) AppID() string {
	return app.appId
}

// Version 实现版本
func (app HadeApp) Version() string {
	return app.configMap["version"]
}

// BaseFolder 表示基础目录，可以代表开发场景的目录，也可以代表运行时候的目录
func (app HadeApp) BaseFolder() string {
	if app.baseFolder != "" {
		return app.baseFolder
	}

	// 如果参数也没有，使用默认的当前路径
	return util.GetExecDirectory()
}

// ConfigFolder  表示配置文件地址
func (app HadeApp) ConfigFolder() string {
	if val, ok := app.configMap["config_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "config")
}

// LogFolder 表示日志存放地址
func (app HadeApp) LogFolder() string {
	if val, ok := app.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(app.StorageFolder(), "log")
}

func (app HadeApp) HttpFolder() string {
	if val, ok := app.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app", "http")
}

func (app HadeApp) ConsoleFolder() string {
	if val, ok := app.configMap["console_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app", "console")
}

func (app HadeApp) StorageFolder() string {
	if val, ok := app.configMap["storage_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "storage")
}

// ProviderFolder 定义业务自己的服务提供者地址
func (app HadeApp) ProviderFolder() string {
	if val, ok := app.configMap["provider_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app", "provider")
}

// MiddlewareFolder 定义业务自己定义的中间件
func (app HadeApp) MiddlewareFolder() string {
	if val, ok := app.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(app.HttpFolder(), "middleware")
}

// CommandFolder 定义业务定义的命令
func (app HadeApp) CommandFolder() string {
	if val, ok := app.configMap["command_folder"]; ok {
		return val
	}
	return filepath.Join(app.ConsoleFolder(), "command")
}

// RuntimeFolder 定义业务的运行中间态信息
func (app HadeApp) RuntimeFolder() string {
	if val, ok := app.configMap["runtime_folder"]; ok {
		return val
	}
	return filepath.Join(app.StorageFolder(), "runtime")
}

// TestFolder 定义测试需要的信息
func (app HadeApp) TestFolder() string {
	if val, ok := app.configMap["test_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "test")
}

// DeployFolder 定义测试需要的信息
func (app HadeApp) DeployFolder() string {
	if val, ok := app.configMap["deploy_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "deploy")
}

// AppFolder 代表app目录
func (app HadeApp) AppFolder() string {
	if val, ok := app.configMap["app_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app")
}

// LoadAppConfig 加载配置map
func (app HadeApp) LoadAppConfig(kv map[string]string) {
	for key, val := range kv {
		app.configMap[key] = val
	}
}

// LoadConfig 加载配置
func (app HadeApp) LoadConfig(val []byte) error {
	return nil
}

// NewHadeApp 初始化HadeApp
func NewHadeApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}

	// 有两个参数，一个是容器，一个是baseFolder
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)
	// 如果没有设置，则使用参数
	if baseFolder == "" {
		baseFolder = util.GetExecDirectory()
	}
	appId := uuid.New().String()
	configMap := map[string]string{}

	return HadeApp{baseFolder: baseFolder, container: container, appId: appId, configMap: configMap}, nil
}
