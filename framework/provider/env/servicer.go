package env

import (
	"bufio"
	"errors"
	"hade/framework/contact"
	"io"
	"os"
	"path"
	"strings"
)

// HadeEnv 是 Env 的具体实现
type HadeEnv struct {
	folder string            // 代表.env所在的目录
	maps   map[string]string // 保存所有的环境变量
}

// NewHadeEnv 创建HadeEnv服务
func NewHadeEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("NewHadeEnv param error")
	}

	folder := params[0].(string)
	// 实例化
	hadeEnv := &HadeEnv{
		folder: folder,
		// 实例化环境变量，APP_ENV默认设置为开发环境
		maps: map[string]string{"APP_ENV": contact.EnvDevelopment},
	}
	// 从配置文件中读取配置
	file := path.Join(folder, ".env")
	fi, err := os.Open(file)
	if err == nil {
		defer fi.Close()

		br := bufio.NewReader(fi)
		for {
			line, _, err := br.ReadLine()
			if err == io.EOF {
				break
			}

			s := strings.SplitN(string(line), "=", 2)
			if len(s) != 2 {
				continue
			}
			hadeEnv.maps[s[0]] = s[1]
		}
	}
	// 获取当前程序的环境变量，并且覆盖.env文件下的变量
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) < 2 {
			continue
		}
		hadeEnv.maps[pair[0]] = pair[1]
	}

	return hadeEnv, nil
}

// AppEnv 获取 APP_ENV 这个环境变量，这个环境变量代表当前应用所在的环境
func (hadeEnv *HadeEnv) AppEnv() string {
	return hadeEnv.maps["APP_ENV"]
}

// IsExist 判断某个环境变量是否存在
func (hadeEnv *HadeEnv) IsExist(name string) bool {
	_, isExist := hadeEnv.maps[name]
	return isExist
}

// Get 获取某个环境变量，如果没有设置，则返回空字符串
func (hadeEnv *HadeEnv) Get(name string) string {
	return hadeEnv.maps[name]
}

// All 获取所有的环境变量
func (hadeEnv *HadeEnv) All() map[string]string {
	return hadeEnv.maps
}
