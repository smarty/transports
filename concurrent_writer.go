package transports

import (
	"io"
	"sync"
)

type ConcurrentWriter struct {
	inner  io.WriteCloser
	mutex  *sync.RWMutex
	closed bool
}

func NewConcurrentWriter(inner io.WriteCloser) io.WriteCloser {
	return &ConcurrentWriter{inner: inner, mutex: &sync.RWMutex{}}
}

func (this *ConcurrentWriter) Write(buffer []byte) (int, error) {
	this.mutex.RLock()
	written, err := this.write(buffer)
	this.mutex.RUnlock()
	return written, err
}
func (this *ConcurrentWriter) write(buffer []byte) (int, error) {
	if this.closed {
		return 0, ErrClosedSocket
	} else {
		return this.inner.Write(buffer)
	}
}

func (this *ConcurrentWriter) Close() error {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.closed {
		return nil
	}
	this.closed = true
	return this.inner.Close()
}
