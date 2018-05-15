package transports

import (
	"log"
	"net"
)

type Dialer interface {
	Dial(string, string) (net.Conn, error)
}

func DefaultDialer() Dialer {
	return &net.Dialer{}
}

////////////////////////////////////////////////////

func DefaultTCPListener(address string) net.Listener {
	if listener, err := NewTCPListener(address); err != nil {
		log.Panic(err)
		return nil
	} else {
		return listener
	}
}

func NewTCPListener(address string) (net.Listener, error) {
	if resolved, err := net.ResolveTCPAddr("tcp", address); err != nil {
		return nil, err
	} else {
		return net.ListenTCP("tcp", resolved)
	}
}
