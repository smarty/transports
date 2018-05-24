package transports

import "crypto/tls"

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
