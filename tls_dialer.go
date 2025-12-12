package transports

import (
	"crypto/tls"
	"log"
	"net"
)

type TLSDialer struct {
	config *tls.Config
	dialer *net.Dialer
}

func NewTLSDialer(dialer *net.Dialer, options ...TLSDialerOption) Dialer {
	this := &TLSDialer{dialer: dialer, config: DefaultTLSConfig()}
	for _, option := range options {
		option(this)
	}
	return *this
}

func (this TLSDialer) Dial(network, address string) (net.Conn, error) {
	conn, err := tls.DialWithDialer(this.dialer, network, address, this.config)
	if err != nil {
		log.Printf("TLSDialer Dial error: [%s]", err.Error())
	} else {
		log.Printf("TLSDialer Dial success - local address: [%s]", conn.LocalAddr().String())
	}
	return conn, err
}

type TLSDialerOption func(*TLSDialer)

func WithTLSConfig(config *tls.Config) TLSDialerOption {
	return func(this *TLSDialer) { this.config = config }
}
func WithTLSServerName(name string) TLSDialerOption {
	return func(this *TLSDialer) { this.config.ServerName = name }
}
