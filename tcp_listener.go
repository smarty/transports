package transports

import (
	"errors"
	"log"
	"net"
	"strings"
)

type TCPListener struct {
	net.Listener
}

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
	} else if listener, err := net.ListenTCP("tcp", resolved); err != nil {
		return nil, err
	} else {
		return TCPListener{Listener: listener}, nil
	}
}

func (this TCPListener) Accept() (net.Conn, error) {
	if socket, err := this.Listener.Accept(); err == nil {
		return socket, nil
	} else if IsClosedError(err) {
		return nil, ErrClosedSocket
	} else {
		return nil, err
	}
}

func IsClosedError(err error) bool {
	return err != nil && strings.Contains(err.Error(), closedNetworkErrorMessage)
}

// https://github.com/golang/go/issues/4373
// https://github.com/golang/go/issues/19252
const closedNetworkErrorMessage = "use of closed network connection"

var ErrClosedSocket = errors.New(closedNetworkErrorMessage)
