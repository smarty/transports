package transports

import (
	"compress/flate"
	"io"
	"net"
)

type GZipConnection struct {
	net.Conn
	reader io.ReadCloser
	writer *flate.Writer
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
	if written, err := this.writer.Write(buffer); err != nil {
		return 0, err
	} else if err = this.writer.Flush(); err != nil {
		return 0, err
	} else {
		return written, nil
	}
}
func (this *GZipConnection) Close() error {
	this.writer.Close()
	this.reader.Close()
	return this.Conn.Close()
}

const (
	NoCompression      = flate.NoCompression
	BestSpeed          = flate.BestSpeed
	BestCompression    = flate.BestCompression
	DefaultCompression = flate.DefaultCompression
)
