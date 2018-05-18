package transports

import (
	"crypto/tls"
	"net"
)

func NewTLSListener(inner net.Listener, config *tls.Config) net.Listener {
	return tls.NewListener(inner, config)
}
