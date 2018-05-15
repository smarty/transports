package transports

import "net"

type Dialer interface {
	Dial(string, string) (net.Conn, error)
}

func NewTCPListener(address string) (net.Listener, error) {
	if resolved, err := net.ResolveTCPAddr("tcp", address); err != nil {
		return nil, err
	} else {
		return net.ListenTCP("tcp", resolved)
	}
}
