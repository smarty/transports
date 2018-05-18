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

func DefaultTLSConfig() *tls.Config {
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
func TLSConfigWithPEM(filename string) (*tls.Config, error) {
	if cert, err := tls.LoadX509KeyPair(filename, filename); err != nil {
		return nil, err
	} else {
		config := DefaultTLSConfig()
		config.Certificates = []tls.Certificate{cert}
		return config, nil
	}
}

type TLSDialerOption func(*TLSDialer)

func WithTLSConfig(config *tls.Config) TLSDialerOption {
	return func(this *TLSDialer) { this.config = config }
}
