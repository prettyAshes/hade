package socket

import (
	"encoding/json"
	"fmt"

	socketio "github.com/googollee/go-socket.io"
	"github.com/spf13/cast"

	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

func InitServer(allowOriginFunc checkOrigin, handle SocketHandle) *socketio.Server {
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	initRoute(server, handle)

	return server
}

func initRoute(server *socketio.Server, handle SocketHandle) {
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")

		handle.ConnectHandle(s)

		url := s.URL()
		room := cast.ToString(url.Query()["room"][0])

		s.Join(room)
		fmt.Println("connected:", s.ID(), " room: ", room)

		broadcastRoomLen(server, room)
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		handle.ErrorHandle(s, e)
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		handle.DisconnectHandle(s, reason)
		room := s.Rooms()[0]
		s.Leave(room)
		s.Close()
		broadcastRoomLen(server, room)

		fmt.Println("closed:", " room: ", room, " reason: ", reason)
	})

	server.OnEvent("/", "msg", func(s socketio.Conn) string {
		return handle.EventHandle(s)
	})
}

func broadcastRoomLen(server *socketio.Server, room string) {
	res := resp{
		ID:  "roomLen",
		Msg: server.RoomLen("/", room),
	}

	result, _ := json.Marshal(res)

	server.BroadcastToRoom("/", room, "reply", string(result))
}
