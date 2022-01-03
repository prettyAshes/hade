package services

import (
	"hade/framework"
	"hade/framework/contact"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type HadeSingleLog struct {
	HadeLog

	folder string
	file   string
	fd     *os.File
}

func NewHadeSingleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contact.LogLevel)
	ctxFielder := params[2].(contact.CtxFielder)
	formatter := params[3].(contact.Formatter)
	logFolder := params[5].(string)

	log := &HadeSingleLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.folder = logFolder
	log.file = "hade.log"

	fd, err := os.OpenFile(filepath.Join(log.folder, log.file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "open log file err")
	}

	log.SetOutput(fd)
	log.c = c

	return log, nil
}
