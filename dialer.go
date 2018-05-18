package transports

import (
	"net"
	"time"
)

type Dialer interface {
	Dial(string, string) (net.Conn, error)
}

func DefaultDialer(options ...DialerOption) Dialer {
	return NetDialer(options...)
}
func NetDialer(options ...DialerOption) *net.Dialer {
	this := &net.Dialer{Timeout: DefaultDialerTimeout}
	for _, option := range options {
		option(this)
	}
	return this
}

type DialerOption func(this *net.Dialer)

func WithDialTimeout(timeout time.Duration) DialerOption {
	return func(this *net.Dialer) { this.Timeout = timeout }
}

const DefaultDialerTimeout = time.Second * 15
