package main

import (
	"fmt"
	"net"

	"github.com/smartystreets/transports"
)

func main() {
	go func() {
		var dialer transports.Dialer = &net.Dialer{}
		//dialer = transports.NewGZipDialer(dialer, gzip.NoCompression)
		clientSocket(dialer.Dial("tcp", "127.0.0.1:8080"))
	}()

	listener, err := transports.NewTCPListener(
		"127.0.0.1:8080",
		serverSocket,
		transports.ListenWithGZip(0))

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Listening...")
	listener.Listen()
}

func clientSocket(socket net.Conn, err error) {
	if socket != nil {
		defer socket.Close()
	}

	if err != nil {
		fmt.Println("ERROR:", err)
	}

	fmt.Println("Writing")
	_, err = socket.Write([]byte("Hello, World!"))
	if err != nil {
		fmt.Println("Write error:", err)
	}
}

func serverSocket(socket net.Conn, err error) {
	if socket != nil {
		defer socket.Close()
	}

	if err != nil {
		fmt.Println("ERROR:", err)
	}

	buffer := make([]byte, 64)
	_, err = socket.Read(buffer)
	if err != nil {
		fmt.Println("Read error:", err)
	} else {
		fmt.Println(string(buffer))
	}
}
