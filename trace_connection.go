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

func NewTraceConnection(inner net.Conn, name string) net.Conn {
	return TraceConnection{Conn: inner, name: name, address: inner.RemoteAddr()}
}
func (this TraceConnection) Read(buffer []byte) (int, error) {
	read, err := this.Conn.Read(buffer)
	if err != nil && err != io.EOF {
		log.Printf("[INFO] Socket read error for [%s] to [%s]. Error: [%s]\n", this.name, this.address, err)
	}
	return read, err
}
func (this TraceConnection) Write(buffer []byte) (int, error) {
	read, err := this.Conn.Write(buffer)
	if err != nil && err != io.EOF {
		log.Printf("[INFO] Socket write error for [%s] to [%s]. Error: [%s]\n", this.name, this.address, err)
	}
	return read, err
}
func (this TraceConnection) Close() error {
	log.Printf("[INFO] Closing socket [%s] to [%s].\n", this.name, this.address)
	return this.Conn.Close()
}
