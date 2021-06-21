package transports

import (
	"crypto/tls"
	"net"
)

func NewTLSListener(inner net.Listener, config *tls.Config) net.Listener {
	if config == nil {
		return inner
	}

	return tls.NewListener(inner, config)
}
