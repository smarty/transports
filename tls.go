package transports

import (
	"crypto/tls"
	"net"
)

type TLSConnection struct {
	net.Conn
}

func NewTLSClient(inner net.Conn, config *tls.Config) (conn *TLSConnection, err error) {
	connection := tls.Client(inner, config)
	if err = connection.Handshake(); err == nil {
		return &TLSConnection{Conn: connection}, nil
	}

	connection.Close()
	return nil, err
}

func NewTLSServer(inner net.Conn, config *tls.Config) (conn *TLSConnection, err error) {
	connection := tls.Server(inner, config)
	if err = connection.Handshake(); err == nil {
		return &TLSConnection{Conn: connection}, nil
	}

	connection.Close()
	return nil, err
}

////////////////////////////////////////////////////

type TLSListener struct {
	net.Listener
	config *tls.Config
}

func NewTLSListener(inner net.Listener, config *tls.Config) *TLSListener {
	return &TLSListener{Listener: inner, config: config}
}

func (this *TLSListener) Accept() (net.Conn, error) {
	if conn, err := this.Listener.Accept(); err != nil {
		return nil, err
	} else {
		return NewTLSServer(conn, this.config)
	}
}

////////////////////////////////////////////////////

type TLSDialer struct {
	Dialer
	config *tls.Config
}

func NewTLSDialer(inner Dialer, config *tls.Config) Dialer {
	return &TLSDialer{Dialer: inner, config: config}
}

func (this *TLSDialer) Dial(network, address string) (net.Conn, error) {
	if conn, err := this.Dialer.Dial(network, address); err != nil {
		return nil, err
	} else {
		return NewTLSClient(conn, this.config)
	}
}
