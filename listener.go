package transports

import (
	"crypto/tls"
	"net"
)

func NewTCPListener(address string) (net.Listener, error) {
	if resolved, err := net.ResolveTCPAddr("tcp", address); err != nil {
		return nil, err
	} else {
		return net.ListenTCP("tcp", resolved)
	}
}

////////////////////////////////////////////////////

type GZipListener struct {
	net.Listener
	compressionLevel int
}

func NewGZipListener(inner net.Listener, compressionLevel int) *GZipListener {
	return &GZipListener{Listener: inner, compressionLevel: compressionLevel}
}

func (this *GZipListener) Accept() (net.Conn, error) {
	if conn, err := this.Listener.Accept(); err != nil {
		return nil, err
	} else {
		return NewGZipConnection(conn, this.compressionLevel)
	}
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
