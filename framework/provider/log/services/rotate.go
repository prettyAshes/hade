package services

import (
	"fmt"
	"hade/framework"
	"hade/framework/contact"
	"hade/framework/util"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type HadeRotateLog struct {
	HadeLog

	// 日志文件存储目录
	folder string
	// 日志文件名
	file string
}

func NewHadeRotateLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contact.LogLevel)
	ctxFielder := params[2].(contact.CtxFielder)
	formatter := params[3].(contact.Formatter)
	logFolder := params[5].(string)

	configServicer := c.MustGetInstance(contact.ConfigKey).(contact.Config)

	// 如果folder不存在，则创建
	if !util.Exists(logFolder) {
		os.MkdirAll(logFolder, os.ModePerm)
	}

	file := "hade.log"

	// 从配置文件获取date_format信息
	dateFormat := "%Y%m%d%H"
	if configServicer.IsExist("log.date_format") {
		dateFormat = configServicer.GetString("log.date_format")
	}

	linkName := rotatelogs.WithLinkName(filepath.Join(logFolder, file))
	options := []rotatelogs.Option{linkName}

	// 从配置文件获取rotate_count信息
	if configServicer.IsExist("log.rotate_count") {
		rotateCount := configServicer.GetInt("log.rotate_count")
		options = append(options, rotatelogs.WithRotationCount(uint(rotateCount)))
	}

	// 从配置文件获取rotate_size信息
	if configServicer.IsExist("log.rotate_size") {
		rotateSize := configServicer.GetInt("log.rotate_size")
		options = append(options, rotatelogs.WithRotationSize(int64(rotateSize)))
	}

	// 从配置文件获取max_age信息
	if configServicer.IsExist("log.max_age") {
		if maxAgeParse, err := time.ParseDuration(configServicer.GetString("log.max_age")); err == nil {
			options = append(options, rotatelogs.WithMaxAge(maxAgeParse))
		}
	}

	// 从配置文件获取rotate_time信息
	if configServicer.IsExist("log.rotate_time") {
		if rotateTimeParse, err := time.ParseDuration(configServicer.GetString("log.rotate_time")); err == nil {
			options = append(options, rotatelogs.WithRotationTime(rotateTimeParse))
		}
	}

	// 设置基础信息
	log := &HadeRotateLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.folder = logFolder
	log.file = file

	w, err := rotatelogs.New(fmt.Sprintf("%s.%s", filepath.Join(log.folder, log.file), dateFormat), options...)
	if err != nil {
		return nil, errors.Wrap(err, "new rotatelogs error")
	}
	log.SetOutput(w)
	log.c = c
	return log, nil
}
