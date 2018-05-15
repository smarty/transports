package transports

import "net"

type Handler func(net.Conn, error)
