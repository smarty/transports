package transports

import (
	"compress/gzip"
	"io"
	"net"
)

type GZip struct {
	net.Conn
	reader io.ReadCloser
	writer io.WriteCloser
}

func NewGZip(inner net.Conn, level int) (*GZip, error) {
	writer, err := gzip.NewWriterLevel(inner, level)
	if err != nil {
		return nil, err
	}

	reader, err := gzip.NewReader(inner)
	if err != nil {
		return nil, err
	}

	return &GZip{Conn: inner, reader: reader, writer: writer}, nil
}

func (this *GZip) Read(buffer []byte) (int, error) {
	return this.reader.Read(buffer)
}
func (this *GZip) Write(buffer []byte) (int, error) {
	return this.writer.Write(buffer)
}
func (this *GZip) Close() error {
	this.writer.Close()
	this.reader.Close()
	return this.Conn.Close()
}
