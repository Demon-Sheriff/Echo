package main

import (
	"syscall"
	"fmt"
)

const (
	AF_INET = 2 
	SOCK_STREAM = 1
	IPPROTO_TCP = 6
)

type SockAddr struct {
	family int
	port int
	ip [4]byte // Store IP in bytes
}

type TCPListener struct {
	fileDescriptor int
	addr SockAddr
}

type Conn struct {
	fileDescriptor int
}

func main() {
	// addr := SockAddr {
	// 	family: AF_INET,
	// 	port : htons(8080)
	// 	ip : 
	// }

}

func NewTCPListener(addr SockAddr) (*TCPListener, error) {
	//Create socket
	fd, err := syscall.Socket(AF_INET, SOCK_STREAM, IPPROTO_TCP)
	if err != nil {
		return nil, fmt.Errorf("Error creating socket: %v", err)
	}

	//Bind the socket to address
	err = syscall.Bind(fd, &addr)
	if err != nil {
		syscall.Close(fd)
		return nil, fmt.Errorf("Error binding socket: %v", err)
	}

	//Listen for incoming connections
	err = syscall.Listen(fd, 10)
	if err != nil {
		syscall.Close(fd)
		return nil, fmt.Errorf("Error listening on socket: %v", err)
	}

	return &TCPListener{
		fileDescriptor: fd,
		addr: addr
	}, nil
}

func (listener *TCPListener) Accept() (Conn, error) {
	fd, sa, err := syscall.Accept(listener.fd)
	if err != nil {
		return nil, fmt.Errorf("Error accepting connection: %v", err)
	}
	
	return Conn{fileDescriptor : fd}, nil
}