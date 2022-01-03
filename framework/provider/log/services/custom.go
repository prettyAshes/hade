package services

import (
	"hade/framework"
	"hade/framework/contact"
	"io"
)

type HadeCustomLog struct {
	HadeLog
}

func NewHadeCustomLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contact.LogLevel)
	ctxFielder := params[2].(contact.CtxFielder)
	formatter := params[3].(contact.Formatter)
	output := params[4].(io.Writer)

	log := &HadeConsoleLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	log.SetOutput(output)
	log.c = c
	return log, nil
}
