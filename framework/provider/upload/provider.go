package upload

import (
	"hade/framework"
	"hade/framework/contact"

	"github.com/spf13/cast"
)

type HadeUploadProvider struct {
	framework.ServiceProvider

	Folder      string
	MaxFileSize int64
}

// Register 注册方法
func (provider *HadeUploadProvider) Register(container framework.Container) framework.NewInstance {
	return NewHadeUpload
}

// Boot 启动调用
func (provider *HadeUploadProvider) Boot(container framework.Container) error {
	configServicer := container.MustGetInstance(contact.ConfigKey).(contact.Config)

	if provider.Folder == "" && configServicer.IsExist("upload.folder") {
		provider.Folder = configServicer.GetString("upload.folder")
	} else {
		panic("upload folder is not exist")
	}

	if provider.MaxFileSize == 0 && configServicer.IsExist("upload.maxFileSize") {
		provider.MaxFileSize = cast.ToInt64(configServicer.Get("upload.maxFileSize"))
	} else {
		panic("maxFileSize is not exist")
	}

	return nil
}

// IsDefer 是否延迟初始化
func (provider *HadeUploadProvider) IsDefer() bool {
	return true
}

// Params 获取初始化参数
func (provider *HadeUploadProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, provider.Folder, provider.MaxFileSize}
}

// Name 获取字符串凭证
func (provider *HadeUploadProvider) Name() string {
	return contact.UploadKey
}
