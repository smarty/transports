package transports

import (
	"errors"
	"net"
	"time"
)

type dialer interface {
	Dial(string) (net.Conn, error)
}

type AutoClient struct {
	address  string
	dialer   dialer
	active   net.Conn
	disposed bool
}

func NewAutoClient(address string, dialer dialer) *AutoClient {
	return &AutoClient{address: address, dialer: dialer}
}

func (this *AutoClient) Read(buffer []byte) (int, error) {
	if conn, err := this.established(); err != nil {
		return 0, err
	} else {
		return conn.Read(buffer)
	}
}
func (this *AutoClient) Write(buffer []byte) (int, error) {
	if conn, err := this.established(); err != nil {
		return 0, err
	} else {
		return conn.Write(buffer)
	}
}
func (this *AutoClient) LocalAddr() net.Addr {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.LocalAddr()
	}
}
func (this *AutoClient) RemoteAddr() net.Addr {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.RemoteAddr()
	}
}
func (this *AutoClient) SetDeadline(instant time.Time) error {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.SetDeadline(instant)
	}
}
func (this *AutoClient) SetReadDeadline(instant time.Time) error {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.SetReadDeadline(instant)
	}
}
func (this *AutoClient) SetWriteDeadline(instant time.Time) error {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.SetWriteDeadline(instant)
	}
}

func (this *AutoClient) established() (net.Conn, error) {
	if this.disposed {
		return nil, DisposedConnection
	} else if this.active != nil {
		return this.active, nil
	} else {
		return this.connect()
	}
}
func (this *AutoClient) connect() (net.Conn, error) {
	active, err := this.dialer.Dial(this.address)
	if err != nil {
		return nil, err
	}
	this.active = active
	return active, nil
}

func (this *AutoClient) Close() error {
	this.close()
	this.disposed = true
	return nil
}
func (this *AutoClient) close() {
	active := this.active
	if active != nil {
		this.active.Close()
	}
	this.active = nil
}

var DisposedConnection = errors.New("connection has been permanently disposed")
