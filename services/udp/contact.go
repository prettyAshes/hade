package udp

import "net"

const UDP_MAXLEN = 65536

type HandleFunc func(data []byte, addr string, netConn *net.UDPConn)
