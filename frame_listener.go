package transports

import "net"

type FrameListener struct {
	net.Listener
}

func NewFrameListener(inner net.Listener) net.Listener {
	return FrameListener{Listener: inner}
}

func (this FrameListener) Accept() (net.Conn, error) {
	if socket, err := this.Listener.Accept(); err == nil {
		return NewFrameConnection(socket), nil
	} else {
		return nil, err
	}
}
