package transports

import (
	"errors"
	"io"
	"sync/atomic"
)

type ChannelWriter struct {
	dialer   Dialer
	address  string
	channel  chan []byte
	closed   uint64
	writer   io.WriteCloser
	buffered [][]byte
}

func NewChannelWriter(dialer Dialer, address string, capacity int) io.WriteCloser {
	this := &ChannelWriter{
		dialer:  dialer,
		address: address,
		channel: make(chan []byte, capacity),
	}
	go this.listen()
	return this
}

func (this *ChannelWriter) Write(buffer []byte) (int, error) {
	if atomic.LoadUint64(&this.closed) > 0 {
		return 0, io.EOF
	}

	select {
	case this.channel <- buffer:
		return len(buffer), nil
	default:
		return 0, ErrChannelFull
	}
}
func (this *ChannelWriter) Close() error {
	if atomic.AddUint64(&this.closed, 1) == 1 {
		close(this.channel)
	}
	return nil
}
func (this *ChannelWriter) isClosed() bool {
	return atomic.LoadUint64(&this.closed) > 0
}

func (this *ChannelWriter) listen() {
	defer this.closeWriter()

	for message := range this.channel {
		this.buffered = append(this.buffered, message)
		if len(this.channel) == 0 {
			this.ensureWrite()
		}
	}
}
func (this *ChannelWriter) ensureWrite() {
	for !this.isClosed() {
		if !this.openWriter() {
			continue
		}

		if this.writeBuffer() {
			break // done
		}

		this.closeWriter()
	}
}
func (this *ChannelWriter) writeBuffer() bool {
	for _, message := range this.buffered {
		if _, err := this.writer.Write(message); err != nil {
			return false
		}
	}

	this.buffered = this.buffered[0:0]
	return true
}
func (this *ChannelWriter) openWriter() bool {
	if this.writer != nil {
		return true
	} else if socket, err := this.dialer.Dial("tcp", this.address); err != nil {
		return false
	} else {
		this.writer = socket
		return true
	}
}
func (this *ChannelWriter) closeWriter() {
	if this.writer != nil {
		this.writer.Close()
		this.writer = nil
	}
}

var ErrChannelFull = errors.New("channel full")
