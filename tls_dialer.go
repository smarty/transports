package transports

import (
	"crypto/tls"
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
	return tls.DialWithDialer(this.dialer, network, address, this.config)
}

type TLSDialerOption func(*TLSDialer)

func WithTLSConfig(config *tls.Config) TLSDialerOption {
	return func(this *TLSDialer) { this.config = config }
}
func WithTLSServerName(name string) TLSDialerOption {
	return func(this *TLSDialer) { this.config.ServerName = name }
}
