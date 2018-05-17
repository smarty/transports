package transports

import (
	"log"
	"net"
)

type TraceConnection struct {
	name string
	net.Conn
}

func NewTraceConnection(name string, inner net.Conn) *TraceConnection {
	log.Printf("[INFO] Socket established for [%s] to [%s]\n", name, inner.RemoteAddr())
	return &TraceConnection{name: name, Conn: inner}
}
func (this *TraceConnection) Read(buffer []byte) (int, error) {
	read, err := this.Conn.Read(buffer)
	if err != nil {
		log.Printf("[INFO] Socket read error for [%s]: [%s]\n", this.name, err)
	}
	return read, err
}
func (this *TraceConnection) Write(buffer []byte) (int, error) {
	read, err := this.Conn.Write(buffer)
	if err != nil {
		log.Printf("[INFO] Socket write error for [%s]: [%s]\n", this.name, err)
	}
	return read, err
}

////////////////////////////////////////////////////

type TraceListener struct {
	name string
	net.Listener
}

func NewTraceListener(name string, inner net.Listener) *TraceListener {
	return &TraceListener{name: name, Listener: inner}
}
func (this *TraceListener) Accept() (net.Conn, error) {
	if socket, err := this.Listener.Accept(); err != nil {
		return nil, err
	} else {
		return NewTraceConnection(this.name, socket), nil
	}
}

////////////////////////////////////////////////////

type TraceDialer struct {
	name string
	Dialer
}

func NewTraceDialer(name string, inner Dialer) *TraceDialer {
	return &TraceDialer{name: name, Dialer: inner}
}
func (this *TraceDialer) Dial(network, address string) (net.Conn, error) {
	if socket, err := this.Dialer.Dial(network, address); err != nil {
		return nil, err
	} else {
		return NewTraceConnection(this.name, socket), nil
	}
}
