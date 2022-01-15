package udp

import (
	"net"
)

// RunServicer 运行UDP服务
func RunServicer(addr string, handle HandleFunc) {
	// bind
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}

	// listen
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	defer udpConn.Close()

	// accept
	buf := make([]byte, UDP_MAXLEN)
	for {
		length, netAddr, err := udpConn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}

		handle(buf[:length], netAddr.String(), udpConn)
	}
}
