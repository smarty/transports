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

func DefaultConfig() *tls.Config {
	return &tls.Config{
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		SessionTicketsDisabled:   true,
		CipherSuites: []uint16{
			tls.TLS_FALLBACK_SCSV,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}
func ConfigWithPEM(filename string) (*tls.Config, error) {
	if cert, err := tls.LoadX509KeyPair(filename, filename); err != nil {
		return nil, err
	} else {
		config := DefaultConfig()
		config.Certificates = []tls.Certificate{cert}
		return config, nil
	}
}
