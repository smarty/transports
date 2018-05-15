package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/smartystreets/transports"
)

const address = "127.0.0.1:8081"

func main() {
	go func() {
		time.Sleep(time.Millisecond) // race condition: let the server start
		clientSocket(newDialer().Dial("tcp", address))
	}()

	fmt.Println("[SERVER] Listening...")
	listener := openListener(address)
	for {
		go serverSocket(listener.Accept())
	}
}

func newDialer() transports.Dialer {
	dialer := transports.NewDialer()
	dialer = transports.NewGZipDialer(dialer, transports.BestCompression)
	return dialer
}

func openListener(address string) net.Listener {
	listener := transports.OpenTCPListener(address)
	listener = transports.NewGZipListener(listener, transports.BestCompression)
	return listener
}

func clientSocket(socket net.Conn, err error) {
	if socket != nil {
		defer socket.Close()
	}

	if err != nil {
		fmt.Println("[CLIENT] ERROR:", err)
		return
	}

	_, err = socket.Write([]byte("Hello, World!"))
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
