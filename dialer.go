package transports

import (
	"crypto/tls"
	"net"
)

type Dialer interface {
	Dial(string, string) (net.Conn, error)
}

////////////////////////////////////////////////////

type GZipDialer struct {
	Dialer
	compression int
}

func NewGZipDialer(inner Dialer, compressionLevel int) Dialer {
	return &GZipDialer{Dialer: inner, compression: compressionLevel}
}

func (this *GZipDialer) Dial(network, address string) (net.Conn, error) {
	if conn, err := this.Dialer.Dial(network, address); err != nil {
		return nil, err
	} else {
		return NewGZipConnection(conn, this.compression)
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
