package config

import (
	"bytes"
	"fmt"
	"hade/framework"
	"hade/framework/contact"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cast"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// HadeConfig  表示hade框架的配置文件服务
type HadeConfig struct {
	c        framework.Container    // 容器
	folder   string                 // 文件夹
	keyDelim string                 // 路径的分隔符，默认为点
	lock     sync.RWMutex           // 配置文件读写锁
	envMaps  map[string]string      // 所有的环境变量
	confMaps map[string]interface{} // 配置文件结构，key为文件名
	confRaws map[string][]byte      // 配置文件的原始信息
}

func NewHadeConfig(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	envFolder := params[1].(string)
	envMaps := params[2].(map[string]string)

	hadeConf := &HadeConfig{
		c:        container,
		folder:   envFolder,
		keyDelim: ".",
		lock:     sync.RWMutex{},
		envMaps:  envMaps,
		confMaps: map[string]interface{}{},
		confRaws: map[string][]byte{},
	}

	if _, err := os.Stat(envFolder); os.IsNotExist(err) {
		return hadeConf, nil
	}

	files, err := ioutil.ReadDir(envFolder)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, file := range files {
		fileName := file.Name()

		err := hadeConf.loadConfigFile(envFolder, fileName)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// 监控文件夹文件
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = watch.Add(envFolder)
	if err != nil {
		return nil, err
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		for {
			select {
			case ev := <-watch.Events:
				{
					//判断事件发生的类型
					// Create 创建
					// Write 写入
					// Remove 删除
					path, _ := filepath.Abs(ev.Name)
					index := strings.LastIndex(path, string(os.PathSeparator))
					folder := path[:index]
					fileName := path[index+1:]

					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("创建文件 : ", ev.Name)
						hadeConf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("写入文件 : ", ev.Name)
						hadeConf.loadConfigFile(folder, fileName)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("删除文件 : ", ev.Name)
						hadeConf.removeConfigFile(folder, fileName)
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error : ", err)
					return
				}
			}
		}
	}()

	return hadeConf, nil
}

func (hadeConfig *HadeConfig) loadConfigFile(envFolder, fileName string) error {
	hadeConfig.lock.Lock()
	defer hadeConfig.lock.Unlock()

	//  判断文件是否以yaml或者yml作为后缀
	s := strings.Split(fileName, ".")
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]

		// 读取文件内容
		bf, err := ioutil.ReadFile(filepath.Join(envFolder, fileName))
		if err != nil {
			return err
		}
		// 直接针对文本做环境变量的替换
		bf = replace(bf, hadeConfig.envMaps)
		// 解析对应的文件
		c := map[string]interface{}{}
		if err := yaml.Unmarshal(bf, &c); err != nil {
			return err
		}
		hadeConfig.confMaps[name] = c
		hadeConfig.confRaws[name] = bf

		// 读取app.path中的信息，更新app对应的folder
		if name == "app" && hadeConfig.c.IsBind(contact.AppKey) {
			if p, ok := c["path"]; ok {
				appService := hadeConfig.c.MustGetInstance(contact.AppKey).(contact.App)
				appService.LoadAppConfig(cast.ToStringMapString(p))
			}
		}
	}

	return nil
}

// removeConfigFile 删除文件的操作
func (conf *HadeConfig) removeConfigFile(folder string, file string) error {
	conf.lock.Lock()
	defer conf.lock.Unlock()
	s := strings.Split(file, ".")
	// 只有yaml或者yml后缀才执行
	if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
		name := s[0]
		// 删除内存中对应的key
		delete(conf.confRaws, name)
		delete(conf.confMaps, name)
	}
	return nil
}

// replace 表示使用环境变量maps替换context中的env(xxx)的环境变量
func replace(content []byte, maps map[string]string) []byte {
	if maps == nil {
		return content
	}

	// 直接使用ReplaceAll替换。这个性能可能不是最优，但是配置文件加载，频率是比较低的，可以接受
	for key, val := range maps {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}

	return content
}

// GetString 获取string类型的config配置
func (conf *HadeConfig) GetString(key string) string {
	return cast.ToString(conf.find(key))
}

// IsExist check setting is exist
func (conf *HadeConfig) IsExist(key string) bool {
	return conf.find(key) != nil
}

// GetInt 获取int类型配置
func (conf *HadeConfig) GetInt(key string) int {
	return cast.ToInt(conf.find(key))
}

// find 获取配置
func (conf *HadeConfig) find(key string) interface{} {
	conf.lock.RLock()
	defer conf.lock.RUnlock()
	return searchMap(conf.confMaps, strings.Split(key, conf.keyDelim))
}

// searchMap 获取制定path的配置
func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		// 判断这个路径是否为1
		if len(path) == 1 {
			return next
		}

		// 判断下一个路径的类型
		switch next := next.(type) {
		case map[interface{}]interface{}:
			// 如果是interface的map，使用cast进行下value转换
			return searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			// 如果是map[string]，直接循环调用
			return searchMap(next, path[1:])
		default:
			// 否则的话，返回nil
			return nil
		}
	}
	return nil
}
