package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"time"

	"github.com/smartystreets/transports"
)

const (
	certificateFile = ""
	serverName      = ""
	address         = "127.0.0.1:8081"
)

var dialAddress = &url.URL{Scheme: "tcp", Host: address}

func main() {
	go func() {
		time.Sleep(time.Millisecond * 250)
		conn := transports.NewAutoConnection(dialAddress, newDialer())
		clientSocket(conn)
	}()

	fmt.Println("[SERVER] Listening...")
	listener := openListener(address)
	for {
		go serverSocket(listener.Accept())
	}
}

func newDialer() transports.Dialer {
	var dialer transports.Dialer
	dialer = transports.NewTLSDialer(&tls.Config{ServerName: serverName})
	dialer = transports.NewGZipDialer(dialer, transports.BestCompression)
	return dialer
}

func openListener(address string) net.Listener {
	listener := transports.DefaultTCPListener(address)
	listener = openTLSListener(listener)
	listener = transports.NewGZipListener(listener, transports.BestCompression)
	return listener
}
func openTLSListener(inner net.Listener) net.Listener {
	cert, _ := tls.LoadX509KeyPair(certificateFile, certificateFile)
	certs := []tls.Certificate{cert}
	return transports.NewTLSListener(inner, &tls.Config{Certificates: certs})
}

func clientSocket(socket net.Conn) {
	if socket != nil {
		defer socket.Close()
	}

	fmt.Println("[CLIENT] Writing bytes")
	written, err := socket.Write([]byte("Hello, World!"))
	if err != nil {
		fmt.Println("[CLIENT] Write error:", err)
	} else {
		fmt.Println("[CLIENT] Bytes written:", written)
	}
}

func serverSocket(socket net.Conn, err error) {
	if socket != nil {
		defer socket.Close()
	}

	if err != nil {
		fmt.Println("[SERVER] ERROR:", err)
		return
	}

	fmt.Println("[SERVER] Socket accepted, connection established...")

	buffer := make([]byte, 64)
	read, err := socket.Read(buffer)
	if read > 0 {
		fmt.Printf("[SERVER] Received from client [%s]\n", string(buffer[:read]))
	}

	if err != nil && err != io.EOF {
		fmt.Println("[SERVER] ERROR:", err)
	}
}
