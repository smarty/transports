package transports

import "net"

type GZipListener struct {
	net.Listener
	compressionLevel int
}

func NewGZipListener(inner net.Listener, compressionLevel int) *GZipListener {
	return &GZipListener{Listener: inner, compressionLevel: compressionLevel}
}

func (this *GZipListener) Accept() (net.Conn, error) {
	if conn, err := this.Listener.Accept(); err != nil {
		return nil, err
	} else {
		return NewGZipConnection(conn, this.compressionLevel)
	}
}
