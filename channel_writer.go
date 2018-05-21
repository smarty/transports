package transports

import (
	"errors"
	"io"
	"sync/atomic"
)

type ChannelWriter struct {
	inner   io.WriteCloser
	channel chan []byte
	closed  int32
}

func NewChannelWriter(inner io.WriteCloser, capacity int) io.WriteCloser {
	this := &ChannelWriter{inner: inner, channel: make(chan []byte, capacity)}
	go this.listen()
	return this
}

func (this *ChannelWriter) Write(buffer []byte) (int, error) {
	select {
	case this.channel <- buffer:
		return len(buffer), nil
	default:
		return 0, ErrBufferFull
	}
}

func (this *ChannelWriter) listen() {
	defer this.inner.Close()

	for buffer := range this.channel {
		this.write(buffer)
	}
}
func (this *ChannelWriter) write(buffer []byte) bool {
	for !this.isClosed() {
		if _, err := this.inner.Write(buffer); err == nil {
			return true
		}
	}
	return false
}

func (this *ChannelWriter) Close() error {
	if atomic.CompareAndSwapInt32(&this.closed, 0, 1) {
		close(this.channel)
	}
	return nil
}
func (this *ChannelWriter) isClosed() bool {
	return atomic.LoadInt32(&this.closed) > 0
}

var ErrBufferFull = errors.New("buffer full")
