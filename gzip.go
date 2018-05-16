package transports

import (
	"compress/flate"
	"io"
	"net"
)

type GZipConnection struct {
	net.Conn
	reader io.ReadCloser
	writer io.WriteCloser
}

func NewGZipConnection(inner net.Conn, level int) (*GZipConnection, error) {
	if writer, err := flate.NewWriter(inner, level); err != nil {
		return nil, err
	} else {
		return &GZipConnection{Conn: inner, reader: flate.NewReader(inner), writer: writer}, nil
	}
}

func (this *GZipConnection) Read(buffer []byte) (int, error) {
	return this.reader.Read(buffer)
}
func (this *GZipConnection) Write(buffer []byte) (int, error) {
	return this.writer.Write(buffer)
}
func (this *GZipConnection) Close() error {
	this.writer.Close()
	this.reader.Close()
	return this.Conn.Close()
}

////////////////////////////////////////////////////

const (
	NoCompression      = flate.NoCompression
	BestSpeed          = flate.BestSpeed
	BestCompression    = flate.BestCompression
	DefaultCompression = flate.DefaultCompression
)

////////////////////////////////////////////////////

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

////////////////////////////////////////////////////

type GZipDialer struct {
	Dialer
	compression int
}

func NewGZipDialer(inner Dialer, options ...GZipDialerOption) Dialer {
	this := &GZipDialer{Dialer: inner, compression: BestCompression}
	for _, option := range options {
		option(this)
	}
	return this
}

func (this *GZipDialer) Dial(network, address string) (net.Conn, error) {
	if conn, err := this.Dialer.Dial(network, address); err != nil {
		return nil, err
	} else {
		return NewGZipConnection(conn, this.compression)
	}
}

////////////////////////////////////////////////////

type GZipDialerOption func(*GZipDialer)

func WithCompressionLevel(level int) GZipDialerOption {
	return func(this *GZipDialer) { this.compression = level }
}
