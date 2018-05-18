package transports

import (
	"log"
	"net"
)

type TraceListener struct {
	net.Listener
	name string
}

func NewTraceListener(inner net.Listener, name string) *TraceListener {
	return &TraceListener{Listener: inner, name: name}
}
func (this *TraceListener) Accept() (net.Conn, error) {
	if socket, err := this.Listener.Accept(); err != nil {
		return nil, err
	} else {
		log.Printf("[INFO] Socket established for [%s] from [%s] to [%s].\n", this.name, socket.RemoteAddr(), socket.LocalAddr())
		return NewTraceConnection(socket, this.name), nil
	}
}
