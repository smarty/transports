package transports

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strings"
)

type FrameConnection struct {
	net.Conn
}

func NewFrameConnection(inner net.Conn) FrameConnection {
	return FrameConnection{Conn: inner}
}

func (this FrameConnection) Write(buffer []byte) (int, error) {
	payloadSize := len(buffer)
	if payloadSize == 0 {
		return 0, nil
	}

	if payloadSize > MaxWriteSize {
		return 0, WriteTooLarge
	}

	if err := binary.Write(this.Conn, byteOrdering, uint16(payloadSize)); err != nil {
		return 0, err
	}

	return this.Conn.Write(buffer)
}
func (this FrameConnection) ReadHeader() (int, error) {
	var length uint16
	if err := binary.Read(this.Conn, byteOrdering, &length); err != nil {
		return 0, err
	} else {
		return int(length), nil
	}
}
func (this FrameConnection) ReadBody(buffer []byte) (int, error) {
	return io.ReadFull(this.Conn, buffer)
}
func (this FrameConnection) Read(buffer []byte) (int, error) {
	if length, err := this.ReadHeader(); err != nil {
		return 0, err
	} else if length > len(buffer) {
		return 0, io.ErrShortBuffer
	} else {
		return this.ReadBody(buffer[0:length])
	}
}

const MaxWriteSize = 64*1024 - 2

var (
	WriteTooLarge = errors.New("buffer to write larger than the max frame size")
	byteOrdering  = binary.LittleEndian
)

////////////////////////////////////////////////////

type FrameListener struct {
	net.Listener
}

func NewFrameListener(inner net.Listener) FrameListener {
	return FrameListener{Listener: inner}
}

func (this FrameListener) Accept() (net.Conn, error) {
	if socket, err := this.Listener.Accept(); err == nil {
		return NewFrameConnection(socket), nil
	} else if strings.Contains(err.Error(), closedAcceptSocketErrorMessage) {
		return nil, io.EOF
	} else {
		return nil, err
	}
}

////////////////////////////////////////////////////

type FrameDialer struct {
	Dialer
}

func NewFrameDialer(inner Dialer) *FrameDialer {
	return &FrameDialer{Dialer: inner}
}

func (this *FrameDialer) Dial(network, address string) (net.Conn, error) {
	if socket, err := this.Dialer.Dial(network, address); err == nil {
		return NewFrameConnection(socket), nil
	} else {
		return nil, err
	}
}

// https://github.com/golang/go/issues/4373
// https://github.com/golang/go/issues/19252
const closedAcceptSocketErrorMessage = "use of closed network connection"
