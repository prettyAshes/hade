package formatter

import "hade/framework/contact"

func Prefix(level contact.LogLevel) string {
	prefix := ""
	switch level {
	case contact.PanicLevel:
		prefix = "[Panic]"
	case contact.FatalLevel:
		prefix = "[Fatal]"
	case contact.ErrorLevel:
		prefix = "[Error]"
	case contact.WarnLevel:
		prefix = "[Warn]"
	case contact.InfoLevel:
		prefix = "[Info]"
	case contact.DebugLevel:
		prefix = "[Debug]"
	case contact.TraceLevel:
		prefix = "[Trace]"
	}
	return prefix
}
