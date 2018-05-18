package transports

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
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
func (this FrameConnection) ReadHeader() (length uint16, err error) {
	if err = binary.Read(this.Conn, byteOrdering, &length); err != nil {
		return 0, err
	} else {
		return length, nil
	}
}
func (this FrameConnection) ReadBody(buffer []byte) (int, error) {
	return io.ReadFull(this.Conn, buffer)
}
func (this FrameConnection) Read(buffer []byte) (int, error) {
	if _, err := this.ReadHeader(); err != nil {
		return 0, err
	} else {
		return this.ReadBody(buffer)
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

func NewFrameListener(inner net.Listener) *FrameListener {
	return &FrameListener{Listener: inner}
}

func (this *FrameListener) Accept() (net.Conn, error) {
	if socket, err := this.Listener.Accept(); err != nil {
		return nil, err
	} else {
		return NewFrameConnection(socket), nil
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
	if socket, err := this.Dialer.Dial(network, address); err != nil {
		return nil, err
	} else {
		return NewFrameConnection(socket), nil
	}
}
