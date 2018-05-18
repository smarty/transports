package transports

import (
	"io"
	"net"
)

type DisposeConnection struct {
	net.Conn
	cleanup func(io.Closer)
}

func NewDisposeConnection(inner net.Conn, cleanup func(closer io.Closer)) net.Conn {
	return DisposeConnection{Conn: inner, cleanup: cleanup}
}

func (this DisposeConnection) Close() error {
	defer this.cleanup(this.Conn)
	return this.Conn.Close()
}
