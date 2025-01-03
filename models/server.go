package models

import (
	"fmt"
	"net"
)

type Server struct {
	Port    int
	Subnet  [4]byte
	clients []net.Conn // ServerMap => for mapping Client with their conn object(through which server has established connection)
}

func (server *Server) startServer() {

	ip := net.IP(server.Subnet[:]).String()
	address := fmt.Sprintf("%s:%d", ip, server.Port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		fmt.Printf("error creating server %v", err)
	}

	//Start listening to clients trying to connect
	for {
		clientConn, err := listener.Accept()

		if err != nil {
			fmt.Printf("error connecting to client %v", err)
		}

		server.addClient(clientConn)
		go server.handleClient(clientConn)
	}
}

func (server *Server) addClient(clientConn net.Conn) {
	server.clients = append(server.clients, clientConn)
}

func (server *Server) handleClient(clientConn net.Conn) {

}
