package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

func broadcastServerUDP(addr string, port int) error {
	// Resolve the broadcast address for the local subnet
	broadcast_id := "192.168.255.255"
	broadcastAddr := fmt.Sprintf("%s:%d", broadcast_id, port) // Broadcasting to all devices in the subnet
	udpAddr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	// laddr, err := net.ResolveUDPAddr("udp","192.168.255.255:6000")
	if err != nil {
		return fmt.Errorf("error resolving UDP address: %w", err)
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("error dialing UDP: %w", err)
	}
	defer conn.Close()

	// Enable broadcast on the socket using syscall
	// This step ensures that the socket is set for broadcasting
	file, err := conn.File()
	if err != nil {
		return fmt.Errorf("error getting file descriptor: %w", err)
	}
	defer file.Close()

	// Set the socket option for broadcast
	err = syscall.SetsockoptInt(int(file.Fd()), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	if err != nil {
		return fmt.Errorf("error setting socket broadcast option: %w", err)
	}

	// Broadcast messages
	for {
		// Send a broadcast message
		message := []byte(fmt.Sprintf("Connect with me at 192.168.3.81:8000", ))
		_, err := conn.Write(message)
		if err != nil {
			return fmt.Errorf("error sending message: %w", err)
		}

		// Print a message to indicate that the broadcast was sent
		fmt.Println("Broadcasting message:", string(message))

		// Wait before sending the next broadcast
		time.Sleep(2 * time.Second) // Broadcast every 2 seconds
	}
}

func main() {
	err := broadcastServerUDP("192.168.255.255", 6000)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
