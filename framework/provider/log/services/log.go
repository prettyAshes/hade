package services

import (
	"context"
	"hade/framework"
	"hade/framework/contact"
	"io"
	"time"

	pkgLog "log"
)

type HadeLog struct {
	// 日志级别
	level contact.LogLevel
	// 日志输出格式方法
	formatter contact.Formatter
	// 日志context上下文信息获取函数
	ctxFielder contact.CtxFielder
	// 日志输出信息
	output io.Writer
	// 容器
	c framework.Container
}

// IsLevelEnable 判断这个级别是否可以打印
func (hadeLog *HadeLog) IsLevelEnable(level contact.LogLevel) bool {
	return level <= hadeLog.level
}

// logf 日志打印的底层函数
func (hadeLog *HadeLog) logf(level contact.LogLevel, ctx context.Context, msg string, fields [][]interface{}) error {
	if !hadeLog.IsLevelEnable(level) {
		return nil
	}

	if hadeLog.ctxFielder != nil {
		ctxFields := hadeLog.ctxFielder(ctx)
		for k, v := range ctxFields {
			fields = append(fields, []interface{}{k, v})
		}
	}

	result, err := hadeLog.formatter(level, time.Now(), msg, fields)
	if err != nil {
		return err
	}

	if level == contact.PanicLevel {
		pkgLog.Panicln(string(result))
		return nil
	}

	hadeLog.output.Write(result)
	hadeLog.output.Write([]byte("\r\n"))

	return nil
}

func (hadeLog *HadeLog) Panic(ctx context.Context, msg string, fields [][]interface{}) {
	hadeLog.logf(contact.PanicLevel, ctx, msg, fields)
}

func (hadeLog *HadeLog) Fatal(ctx context.Context, msg string, fields [][]interface{}) {
	hadeLog.logf(contact.FatalLevel, ctx, msg, fields)
}

func (hadeLog *HadeLog) Error(ctx context.Context, msg string, fields [][]interface{}) {
	hadeLog.logf(contact.ErrorLevel, ctx, msg, fields)
}

func (hadeLog *HadeLog) Warn(ctx context.Context, msg string, fields [][]interface{}) {
	hadeLog.logf(contact.WarnLevel, ctx, msg, fields)
}

func (hadeLog *HadeLog) Info(ctx context.Context, msg string, fields [][]interface{}) {
	hadeLog.logf(contact.InfoLevel, ctx, msg, fields)
}

func (hadeLog *HadeLog) Debug(ctx context.Context, msg string, fields [][]interface{}) {
	hadeLog.logf(contact.DebugLevel, ctx, msg, fields)
}

func (hadeLog *HadeLog) Trace(ctx context.Context, msg string, fields [][]interface{}) {
	hadeLog.logf(contact.TraceLevel, ctx, msg, fields)
}

func (hadeLog *HadeLog) SetLevel(level contact.LogLevel) {
	hadeLog.level = level
}

func (hadeLog *HadeLog) SetCtxFielder(handler contact.CtxFielder) {
	hadeLog.ctxFielder = handler
}

func (hadeLog *HadeLog) SetFormatter(formatter contact.Formatter) {
	hadeLog.formatter = formatter
}

func (hadeLog *HadeLog) SetOutput(out io.Writer) {
	hadeLog.output = out
}

func (hadeLog *HadeLog) GetOutPut() io.Writer {
	return hadeLog.output
}
