package transports

import (
	"crypto/tls"
	"net"
)

func NewTLSListener(inner net.Listener, config *tls.Config) net.Listener {
	return tls.NewListener(inner, config)
}

////////////////////////////////////////////////////

type TLSDialer struct {
	config *tls.Config
}

func NewTLSDialer(config *tls.Config) Dialer {
	return &TLSDialer{config: config}
}

func (this *TLSDialer) Dial(network, address string) (net.Conn, error) {
	return tls.Dial(network, address, this.config)
}
