package socket

import (
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

// checkOrigin 检查origin
type checkOrigin func(r *http.Request) bool

// resp socket响应结构
type resp struct {
	ID  string      `json:"id"`
	Msg interface{} `json:"msg"`
}

// socketHandle socket执行函数
type SocketHandle interface {
	ConnectHandle(s socketio.Conn)
	ErrorHandle(s socketio.Conn, e error)
	DisconnectHandle(s socketio.Conn, reason string)
	EventHandle(s socketio.Conn) string
}
