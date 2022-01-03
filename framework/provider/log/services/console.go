package services

import (
	"hade/framework"
	"hade/framework/contact"
	"os"
)

type HadeConsoleLog struct {
	HadeLog
}

func NewHadeConsoleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contact.LogLevel)
	ctxFielder := params[2].(contact.CtxFielder)
	formatter := params[3].(contact.Formatter)

	log := &HadeLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	log.c = c

	return log, nil
}
