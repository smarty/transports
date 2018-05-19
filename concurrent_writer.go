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
	if this.closed {
		this.mutex.RUnlock()
		return 0, ErrClosedSocket
	}

	written, err := this.inner.Write(buffer)
	this.mutex.RUnlock()
	return written, err
}

func (this *ConcurrentWriter) Close() (err error) {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return
	}

	this.closed = true
	err = this.inner.Close()
	this.mutex.Unlock()
	return err
}
