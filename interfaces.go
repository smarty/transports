package transports

import "net"

type Dialer interface {
	Dial(string, string) (net.Conn, error)
}
