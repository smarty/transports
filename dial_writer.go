package transports

import "io"

type DialWriter struct {
	dialer  Dialer
	address string
	writer  io.WriteCloser
}

func NewDialWriter(dialer Dialer, address string) *DialWriter {
	return &DialWriter{dialer: dialer, address: address}
}

func (this *DialWriter) Write(buffer []byte) (int, error) {
	if this.writer != nil {
		return this.write(buffer)
	} else if writer, err := this.dialer.Dial("tcp", this.address); err != nil {
		return 0, err
	} else {
		this.writer = writer
		return this.write(buffer)
	}
}
func (this *DialWriter) write(buffer []byte) (int, error) {
	if written, err := this.writer.Write(buffer); err == nil {
		return written, nil
	} else {
		this.writer.Close()
		this.writer = nil
		return written, err
	}
}

func (this *DialWriter) Close() (err error) {
	if this.writer != nil {
		err = this.writer.Close()
		this.writer = nil
	}
	return err
}
