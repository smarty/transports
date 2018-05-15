package transports

import (
	"crypto/tls"
	"net"
	"sync/atomic"
)

type Listener struct {
	state    uint64
	inner    net.Listener
	callback Handler
}

func NewTCPListener(address string, callback Handler, options ...ListenerOption) (*Listener, error) {
	if resolved, err := net.ResolveTCPAddr("tcp", address); err != nil {
		return nil, err
	} else if listener, err := net.ListenTCP("tcp", resolved); err != nil {
		return nil, err
	} else {
		return NewListener(listener, callback, options...), nil
	}
}

func NewListener(inner net.Listener, callback Handler, options ...ListenerOption) *Listener {
	this := &Listener{inner: inner, callback: callback}

	for _, option := range options {
		option(this)
	}

	return this
}

func (this *Listener) Listen() {
	for this.isOpen() {
		if socket, err := this.inner.Accept(); err == nil {
			go this.callback(socket, nil)
		}
	}
}

func (this *Listener) Close() error {
	err := this.inner.Close()
	atomic.StoreUint64(&this.state, 1)
	return err
}
func (this *Listener) isOpen() bool {
	return atomic.LoadUint64(&this.state) == 0
}

////////////////////////////////////////////////////

type Handler func(net.Conn, error)
type ListenerOption func(this *Listener)

func WithTLSServer(config *tls.Config) ListenerOption {
	return func(this *Listener) {
		callback := this.callback
		this.callback = func(socket net.Conn, err error) {
			if err == nil {
				socket, err = NewTLSServer(socket, config)
			}
			callback(nil, err)
		}
	}
}
func WithTLSClient(config *tls.Config) ListenerOption {
	return func(this *Listener) {
		callback := this.callback
		this.callback = func(socket net.Conn, err error) {
			if err == nil {
				socket, err = NewTLSClient(socket, config)
			}
			callback(nil, err)
		}
	}
}
func WithGZip(level int) ListenerOption {
	return func(this *Listener) {
		callback := this.callback
		this.callback = func(socket net.Conn, err error) {
			if err == nil {
				socket, err = NewGZip(socket, level)
			}
			callback(nil, err)
		}
	}
}
