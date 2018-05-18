package transports

import "net"

type FrameDialer struct {
	Dialer
}

func NewFrameDialer(inner Dialer) Dialer {
	return FrameDialer{Dialer: inner}
}

func (this FrameDialer) Dial(network, address string) (net.Conn, error) {
	if socket, err := this.Dialer.Dial(network, address); err == nil {
		return NewFrameConnection(socket), nil
	} else {
		return nil, err
	}
}
