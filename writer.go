package transports

import (
	"context"
	"crypto/tls"
	"io"
	"net"
)

type dialWriter struct {
	net.Conn
	dialContext context.Context
	dialNetwork string
	dialAddress string
	dialer      Dialer
	tlsConfig   *tls.Config
	logger      Logger
}

func newDialWriter(config writerConfiguration) io.WriteCloser {
	return &dialWriter{
		dialContext: config.DialContext,
		dialNetwork: config.DialNetwork,
		dialAddress: config.DialAddress,
		dialer:      config.Dialer,
		logger:      config.Logger,
	}
}

func (this *dialWriter) Write(buffer []byte) (int, error) {
	if this.Conn == nil {
		if conn, err := this.dialer.DialContext(this.dialContext, this.dialNetwork, this.dialAddress); err != nil {
			return 0, err
		} else if this.tlsConfig != nil {
			this.Conn = tls.Client(conn, this.tlsConfig)
		} else {
			this.Conn = conn
		}
	}

	if written, err := this.Conn.Write(buffer); err == nil {
		return written, nil
	} else {
		_ = this.Conn.Close()
		this.Conn = nil
		return written, err
	}
}
func (this *dialWriter) Close() error {
	if this.Conn == nil {
		return nil
	}

	defer func() { this.Conn = nil }()
	return this.Conn.Close()
}
