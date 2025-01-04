package models

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

type ClientMap struct {
	mu sync.RWMutex
	clients map[string] *Client
} 

func (cM *ClientMap) AddClient(room_key string, client *Client) {

	cM.mu.Lock()
	defer cM.mu.Unlock()

	if _, exists := cM.clients[room_key]; exists {
		fmt.Printf("room with key %v already exists\n", room_key)
		return
	}

	cM.clients[room_key] = client
	fmt.Printf("connected with room having room_key = %v\n", room_key)
}

func (cM *ClientMap) GetClient(room_key string) bool {

	cM.mu.RLock()
	defer cM.mu.RUnlock()

	_, exists := cM.clients[room_key]
	return exists;
}	

func (cM *ClientMap) RemoveClient(room_key string) {

	cM.mu.Lock()
	defer cM.mu.Unlock()

	if _, exists := cM.clients[room_key]; exists {
		delete(cM.clients, room_key)
		fmt.Printf("Removed the client with room_key = %v\n", room_key)
		return
	}
	fmt.Printf("The chat-room does not exist\n")
}


type Server struct {
	Port int
	Subnet [4]byte
	// ServerMap 
}

// write message to a single client 
func (server *Server) SendMessage(conn net.Conn, done chan bool) {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("send message to the client : ");

		if !scanner.Scan() {
			fmt.Println("Input error or EOF")
			done <- true
			return
		}

		message := scanner.Text() + "\n"


		// send the message to the client
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending the data to the client", err)
			done <- true
			return 
		}
	}
}

func (server *Server) RecvMessage(conn net.Conn, done chan bool) {

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading the input from the client");
			done <- true
			return 
		}

		fmt.Println("client : ", message + "\n")
	}

}