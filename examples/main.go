package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"

	"github.com/smartystreets/transports"
)

func main() {
	go func() {
		var dialer transports.Dialer = &net.Dialer{}
		dialer = transports.NewTLSDialer(dialer, &tls.Config{InsecureSkipVerify: true})
		//dialer = transports.NewGZipDialer(dialer, transports.BestCompression)
		clientSocket(dialer.Dial("tcp", "127.0.0.1:8080"))
	}()

	cert, _ := tls.LoadX509KeyPair("", "")
	listener, err := transports.NewTCPListener(
		"127.0.0.1:8080",
		serverSocket,
		//transports.ListenWithGZip(0),
		transports.ListenWithTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
		}))

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
		return
	}

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
		return
	}

	buffer := make([]byte, 64)
	read, err := socket.Read(buffer)
	if read > 0 {
		fmt.Println(string(buffer[:read]))
	}

	if err != nil && err != io.EOF {
		fmt.Println("ERROR:", err)
	}
}
