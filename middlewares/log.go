package middlewares

import (
	"bytes"
	"hade/framework"
	"hade/framework/contact"
	"hade/framework/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger(container framework.Container) gin.HandlerFunc {
	logService := container.MustGetInstance(contact.LogKey).(contact.Log)
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		logService.Info(c, "response", [][]interface{}{
			{"url", c.Request.URL},
			{"body", blw.body.String()},
		})
	}
}
