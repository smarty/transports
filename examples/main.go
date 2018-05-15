package main

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"time"

	"github.com/smartystreets/transports"
)

const address = "127.0.0.1:8081"

var dialAddress = &url.URL{Scheme: "tcp", Host: address}

func main() {
	go func() {
		conn := transports.NewAutoConnection(dialAddress, newDialer())
		clientSocket(conn)
	}()

	time.Sleep(time.Second)

	fmt.Println("[SERVER] Listening...")
	listener := openListener(address)
	for {
		go serverSocket(listener.Accept())
	}
}

func newDialer() transports.Dialer {
	dialer := transports.DefaultDialer()
	dialer = transports.NewGZipDialer(dialer, transports.BestCompression)
	return dialer
}

func openListener(address string) net.Listener {
	listener := transports.DefaultTCPListener(address)
	listener = transports.NewGZipListener(listener, transports.BestCompression)
	return listener
}

func clientSocket(socket net.Conn) {
	if socket != nil {
		defer socket.Close()
	}

	_, err := socket.Write([]byte("Hello, World!"))
	if err != nil {
		fmt.Println("[CLIENT] Write error:", err)
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
