package transports

import (
	"errors"
	"net"
	"net/url"
	"time"
)

type AutoConnection struct {
	address  *url.URL
	dialer   Dialer
	active   net.Conn
	disposed bool
}

func NewAutoConnection(address *url.URL, dialer Dialer) *AutoConnection {
	return &AutoConnection{address: address, dialer: dialer}
}

func (this *AutoConnection) Read(buffer []byte) (int, error) {
	if conn, err := this.established(); err != nil {
		return 0, err
	} else {
		return conn.Read(buffer)
	}
}
func (this *AutoConnection) Write(buffer []byte) (int, error) {
	if conn, err := this.established(); err != nil {
		return 0, err
	} else {
		return conn.Write(buffer)
	}
}
func (this *AutoConnection) LocalAddr() net.Addr {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.LocalAddr()
	}
}
func (this *AutoConnection) RemoteAddr() net.Addr {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.RemoteAddr()
	}
}
func (this *AutoConnection) SetDeadline(instant time.Time) error {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.SetDeadline(instant)
	}
}
func (this *AutoConnection) SetReadDeadline(instant time.Time) error {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.SetReadDeadline(instant)
	}
}
func (this *AutoConnection) SetWriteDeadline(instant time.Time) error {
	if conn, err := this.established(); err != nil {
		return nil
	} else {
		return conn.SetWriteDeadline(instant)
	}
}

func (this *AutoConnection) established() (net.Conn, error) {
	if this.disposed {
		return nil, DisposedConnection
	} else if this.active != nil {
		return this.active, nil
	} else {
		return this.connect()
	}
}
func (this *AutoConnection) connect() (net.Conn, error) {
	active, err := this.dialer.Dial(this.address.Scheme, this.address.Host)
	if err != nil {
		return nil, err
	}
	this.active = active
	return active, nil
}

func (this *AutoConnection) Close() error {
	this.close()
	this.disposed = true
	return nil
}
func (this *AutoConnection) close() {
	active := this.active
	if active != nil {
		this.active.Close()
	}
	this.active = nil
}

var DisposedConnection = errors.New("connection has been permanently disposed")
