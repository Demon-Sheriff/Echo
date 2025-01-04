package models

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
)

type ClientMap struct {
	mu      sync.RWMutex
	clients map[*Client]net.Conn
}

func (cM *ClientMap) AddClient(conn net.Conn, client *Client) {

	cM.mu.Lock()
	defer cM.mu.Unlock()

	if _, exists := cM.clients[client]; exists {
		fmt.Printf("client with conn %v already exists\n", conn)
		return
	}

	cM.clients[client] = conn
	fmt.Printf("connected with client having conn = %v\n", conn)
}

func (cM *ClientMap) GetClient(client *Client) bool {

	cM.mu.RLock()
	defer cM.mu.RUnlock()

	_, exists := cM.clients[client]
	return exists
}

func (cM *ClientMap) RemoveClient(client *Client) {

	cM.mu.Lock()
	defer cM.mu.Unlock()

	if _, exists := cM.clients[client]; exists {
		delete(cM.clients, client)
		fmt.Printf("Removed the client %v\n", client)
		return
	}
	fmt.Printf("The chat-room does not exist\n")
}

type Server struct {
	Port   int
	Subnet [4]byte
	ClientMap
}

// write message to a single client
func (server *Server) SendMessage(conn net.Conn, done chan bool) {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("send message to the client : ")

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

	decoder := json.NewDecoder(conn)
	var initMessage Message

	if err := decoder.Decode(&initMessage); err != nil {
		fmt.Println("Error reading initial username message:", err)
		return
	}

	// declare new client
	client := Client{
		Client_name: initMessage.Text,
	}

	// add client in clients map
	server.ClientMap.AddClient(conn, &client)
}

func (server *Server) start(chat ChatInterface) {

	ip := net.IP(server.Subnet[:]).String()
	address := fmt.Sprintf("%s:%d", ip, server.Port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		fmt.Printf("error creating server %v", err)
	}

	// define channels for communicating between clients goroutines
	messageChannel := make(chan Message)
	clientChannel := make(chan Client)
	removeClientChannel := make(chan Client)

	// manage channels and update Chat Interface accordinly
	go manageClients(messageChannel, clientChannel, removeClientChannel, chat)

	// start listening to clients
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf("error connecting to client %v", err)
		}

		done := make(chan bool)

		// get the first message to get Client's info
		server.RecvMessage(conn, done)

		// handle each client in separate goroutine
		go handleClient(conn, messageChannel, removeClientChannel)
	}
}

func handleClient(conn net.Conn, messageChannel chan Message, removeClientChannel chan Client) {

}

func manageClients(messageChannel chan Message, clientChannel chan Client, removeClientChannel chan Client, chat ChatInterface) {

}
