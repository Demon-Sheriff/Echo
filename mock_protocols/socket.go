package main

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
)

type Message struct {
	Sender  int
	Message string
}

func main() {

	addr := syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{127, 0, 0, 1},
	}

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Println("Error creating socket:", err)
		return
	}
	defer syscall.Close(fd)
	fmt.Println("Socket created with file descriptor:", fd)

	err = syscall.Bind(fd, &addr)
	if err != nil {
		fmt.Println("Error binding:", err)
		return
	}

	err = syscall.Listen(fd, 5)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	fmt.Println("Socket is listening for connections:")

	messageChannel := make(chan Message)
	clientChannel := make(chan int)
	removeClientChannel := make(chan int)

	//manage channel messages (connects, disconnects, client-messages)
	go manageClients(messageChannel, clientChannel, removeClientChannel)

	//input reading from server
	go serverReader(messageChannel, fd)

	for {
		clientFd, _, err := syscall.Accept(fd)
		if err != nil {
			fmt.Println("Error accepting connection", err)
			return
		}

		//send client channel with new connection message
		clientChannel <- clientFd
		go handleClient(clientFd, messageChannel, removeClientChannel)
	}
}

func manageClients(messageChannel chan Message, clientChannel chan int, removeClientChannel chan int) {
	messageHistory := []string{}
	clients := make(map[int]bool)

	for {
		select {
		case clientFd := <-clientChannel:
			clients[clientFd] = true
			for _, msg := range messageHistory {
				syscall.Write(clientFd, []byte(msg))
			}

		case clientFd := <-removeClientChannel:
			delete(clients, clientFd)

		case msg := <-messageChannel:
			messageHistory = append(messageHistory, msg.Message)

			for clientFd := range clients {
				_, err := syscall.Write(clientFd, []byte(msg.Message))
				if err != nil {
					fmt.Println("Error sending message to client:", err)
					//Probably client is disconnected so better to remove it here so that it doesnt cause more errors
					delete(clients, clientFd)
				}
			}
		}
	}
}

func serverReader(messageChannel chan Message, fd int) {
	reader := bufio.NewReader(os.Stdin)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input", err)
			return
		}

		message = message[:len(message)-1]

		if message == "Shut Down" {
			fmt.Println("Server shutting down...")
			os.Exit(0)
		}

		messageChannel <- Message{Sender: fd, Message: "Server:" + message}
	}
}

func handleClient(clientFd int, messageChannel chan Message, removeClientChannel chan int) {
	defer func() {
		removeClientChannel <- clientFd
		syscall.Close(clientFd)
	}()

	for {
		buffer := make([]byte, 1024)
		n, err := syscall.Read(clientFd, buffer)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}

		if n == 0 {
			fmt.Println("Client disconnected:", clientFd)
			return
		}

		msg := string(buffer[:n])
		messageChannel <- Message{Sender: clientFd, Message: fmt.Sprintf("Client %d: %s", clientFd, msg)}
	}
}
