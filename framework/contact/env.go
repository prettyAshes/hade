package contact

const (
	// EnvProduction 代表生产环境
	EnvProduction = "production"
	// EnvTesting 代表测试环境
	EnvTesting = "testing"
	// EnvDevelopment 代表开发环境
	EnvDevelopment = "development"

	// EnvKey 是环境变量服务字符串凭证
	EnvKey = "hade:env"
)

type Env interface {
	// AppEnv 获取 APP_ENV 这个环境变量，这个环境变量代表当前应用所在的环境
	AppEnv() string
	// IsExist 判断某个环境变量是否存在
	IsExist(name string) bool
	// Get 获取某个环境变量，如果没有设置，则返回空字符串
	Get(name string) string
	// All 获取所有的环境变量
	All() map[string]string
}
