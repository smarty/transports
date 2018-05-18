package transports

import (
	"log"
	"net"
)

type TraceDialer struct {
	Dialer
	name string
}

func NewTraceDialer(inner Dialer, name string) Dialer {
	return TraceDialer{Dialer: inner, name: name}
}
func (this TraceDialer) Dial(network, address string) (net.Conn, error) {
	if socket, err := this.Dialer.Dial(network, address); err != nil {
		return nil, err
	} else {
		log.Printf("[INFO] Socket established for [%s] from [%s] to [%s].\n", this.name, socket.LocalAddr(), socket.RemoteAddr())
		return NewTraceConnection(socket, this.name), nil
	}
}
