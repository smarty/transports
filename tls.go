package transports

import (
	"crypto/tls"
	"net"
)

type TLS struct {
	net.Conn
}

func NewTLSClient(inner net.Conn, config *tls.Config) (conn *TLS, err error) {
	connection := tls.Client(inner, config)
	if err = connection.Handshake(); err == nil {
		return &TLS{Conn: connection}, nil
	}

	connection.Close()
	return nil, err
}

func NewTLSServer(inner net.Conn, config *tls.Config) (conn *TLS, err error) {
	connection := tls.Server(inner, config)
	if err = connection.Handshake(); err == nil {
		return &TLS{Conn: connection}, nil
	}

	connection.Close()
	return nil, err
}
