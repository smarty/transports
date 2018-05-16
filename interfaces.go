package transports

import (
	"log"
	"net"
	"time"
)

type Dialer interface {
	Dial(string, string) (net.Conn, error)
}

func DefaultDialer(options ...DialerOption) Dialer {
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

////////////////////////////////////////////////////

func DefaultTCPListener(address string) net.Listener {
	if listener, err := NewTCPListener(address); err != nil {
		log.Panic(err)
		return nil
	} else {
		return listener
	}
}

func NewTCPListener(address string) (net.Listener, error) {
	if resolved, err := net.ResolveTCPAddr("tcp", address); err != nil {
		return nil, err
	} else {
		return net.ListenTCP("tcp", resolved)
	}
}
