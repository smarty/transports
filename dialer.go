package transports

import (
	"crypto/tls"
	"net"
)

type Dialer struct {
	network  string
	inner    net.Dialer
	callback Handler
}

func NewTCPDialer(callback Handler, options ...DialerOption) *Dialer {
	return NewDialer(net.Dialer{}, callback, options...)
}
func NewDialer(inner net.Dialer, callback Handler, options ...DialerOption) *Dialer {
	this := &Dialer{network: "tcp", inner: inner, callback: callback}

	for _, option := range options {
		option(this)
	}

	return this
}

func (this *Dialer) Dial(address string) (net.Conn, error) {
	return this.inner.Dial(this.network, address)
}

////////////////////////////////////////////////////

type DialerOption func(this *Dialer)

func DialWithTLS(config *tls.Config) DialerOption {
	return func(this *Dialer) {
		callback := this.callback
		this.callback = func(socket net.Conn, err error) {
			if err == nil {
				socket, err = NewTLSClient(socket, config)
			}
			callback(nil, err)
		}
	}
}
func DialWithGZip(level int) DialerOption {
	return func(this *Dialer) {
		callback := this.callback
		this.callback = func(socket net.Conn, err error) {
			if err == nil {
				socket, err = NewGZip(socket, level)
			}
			callback(nil, err)
		}
	}
}
