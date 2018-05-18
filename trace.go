package transports

import (
	"io"
	"log"
	"net"
)

type TraceConnection struct {
	net.Conn
	name    string
	address interface{}
}

func NewTraceConnection(inner net.Conn, name string) *TraceConnection {
	return &TraceConnection{Conn: inner, name: name, address: inner.RemoteAddr()}
}
func (this *TraceConnection) Read(buffer []byte) (int, error) {
	read, err := this.Conn.Read(buffer)
	if err != nil && err != io.EOF {
		log.Printf("[INFO] Socket read error for [%s] to [%s]. Error: [%s]\n", this.name, this.address, err)
	}
	return read, err
}
func (this *TraceConnection) Write(buffer []byte) (int, error) {
	read, err := this.Conn.Write(buffer)
	if err != nil && err != io.EOF {
		log.Printf("[INFO] Socket write error for [%s] to [%s]. Error: [%s]\n", this.name, this.address, err)
	}
	return read, err
}
func (this *TraceConnection) Close() error {
	log.Printf("[INFO] Closing socket [%s] to [%s].\n", this.name, this.address)
	return this.Conn.Close()
}

////////////////////////////////////////////////////

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

////////////////////////////////////////////////////

type TraceDialer struct {
	Dialer
	name string
}

func NewTraceDialer(inner Dialer, name string) *TraceDialer {
	return &TraceDialer{Dialer: inner, name: name}
}
func (this *TraceDialer) Dial(network, address string) (net.Conn, error) {
	if socket, err := this.Dialer.Dial(network, address); err != nil {
		return nil, err
	} else {
		log.Printf("[INFO] Socket established for [%s] from [%s] to [%s].\n", this.name, socket.LocalAddr(), socket.RemoteAddr())
		return NewTraceConnection(socket, this.name), nil
	}
}
