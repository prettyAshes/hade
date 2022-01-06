package formatter

import (
	"bytes"
	"hade/framework/contact"
	"strings"
	"time"

	"github.com/spf13/cast"
)

func JsonFormatter(level contact.LogLevel, t time.Time, msg string, fields [][]interface{}) ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	fields = append(fields,
		[]interface{}{"msg", msg},
		[]interface{}{"level", level},
		[]interface{}{"timestamp", t.Format(time.RFC3339)})

	var l []string
	for _, field := range fields {
		if len(field) != 2 {
			continue
		}
		k := cast.ToString(field[0])
		v := cast.ToString(field[1])
		if k == "" || v == "" {
			continue
		}
		l = append(l, k+": "+v)
	}

	bf.Write([]byte(strings.Join(l, " | ")))
	return bf.Bytes(), nil
}
