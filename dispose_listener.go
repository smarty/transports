package transports

import (
	"io"
	"net"
	"sync"
)

type DisposeListener struct {
	net.Listener
	mutex   *sync.Mutex
	tracked map[io.Closer]struct{}
}

func NewDisposeListener(inner net.Listener) net.Listener {
	return &DisposeListener{
		Listener: inner,
		mutex:    &sync.Mutex{},
		tracked:  make(map[io.Closer]struct{}),
	}
}

func (this *DisposeListener) Accept() (net.Conn, error) {
	if actual, err := this.Listener.Accept(); err != nil {
		return nil, err
	} else {
		this.add(actual)
		return NewDisposeConnection(actual, this.remove), nil
	}
}
func (this *DisposeListener) add(actual net.Conn) {
	this.mutex.Lock()
	this.tracked[actual] = struct{}{}
	this.mutex.Unlock()

}
func (this *DisposeListener) remove(actual io.Closer) {
	this.mutex.Lock()
	delete(this.tracked, actual)
	this.mutex.Unlock()
}

func (this *DisposeListener) Close() error {
	err := this.Listener.Close()

	this.mutex.Lock()
	defer this.mutex.Unlock()

	for actual := range this.tracked {
		actual.Close()
	}

	return err
}
