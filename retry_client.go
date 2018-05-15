package transports

import (
	"errors"
	"net"
	"time"
)

type dialer interface {
	Dial(string) (net.Conn, error)
}

type RetryClient struct {
	address  string
	dialer   dialer
	active   net.Conn
	disposed bool
}

func NewRetryClient(address string, dialer dialer) *RetryClient {
	return &RetryClient{address: address, dialer: dialer}
}

func (this *RetryClient) Read(buffer []byte) (int, error) {
	if conn, err := this.established(); err != nil {
		return 0, err
	} else {
		return conn.Read(buffer)
	}
}
func (this *RetryClient) Write(buffer []byte) (int, error) {
	if conn, err := this.established(); err != nil {
		return 0, err
	} else {
		return conn.Write(buffer)
	}
}
func (this *RetryClient) LocalAddr() net.Addr {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.LocalAddr()
	}
}
func (this *RetryClient) RemoteAddr() net.Addr {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.RemoteAddr()
	}
}
func (this *RetryClient) SetDeadline(instant time.Time) error {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.SetDeadline(instant)
	}
}
func (this *RetryClient) SetReadDeadline(instant time.Time) error {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.SetReadDeadline(instant)
	}
}
func (this *RetryClient) SetWriteDeadline(instant time.Time) error {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.SetWriteDeadline(instant)
	}
}

func (this *RetryClient) established() (net.Conn, error) {
	if this.disposed {
		return nil, DisposedConnection
	} else if this.active != nil {
		return this.active, nil
	} else {
		return this.connect()
	}
}
func (this *RetryClient) connect() (net.Conn, error) {
	active, err := this.dialer.Dial(this.address)
	if err != nil {
		return nil, err
	}
	this.active = active
	return active, nil
}

func (this *RetryClient) Close() error {
	this.close()
	this.disposed = true
	return nil
}
func (this *RetryClient) close() {
	active := this.active
	if active != nil {
		this.active.Close()
	}
	this.active = nil
}

var DisposedConnection = errors.New("connection has been permanently disposed")
