package transports

import "net"

type GZipDialer struct {
	Dialer
	compression int
}

func NewGZipDialer(inner Dialer, options ...GZipDialerOption) Dialer {
	this := &GZipDialer{Dialer: inner, compression: BestCompression}
	for _, option := range options {
		option(this)
	}
	return this
}

func (this *GZipDialer) Dial(network, address string) (net.Conn, error) {
	if conn, err := this.Dialer.Dial(network, address); err != nil {
		return nil, err
	} else {
		return NewGZipConnection(conn, this.compression)
	}
}

type GZipDialerOption func(*GZipDialer)

func WithCompressionLevel(level int) GZipDialerOption {
	return func(this *GZipDialer) { this.compression = level }
}
