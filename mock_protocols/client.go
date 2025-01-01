package main

import (
	"fmt"
	"net"
	"os"
)

func listenForBroadcast(port int) error {
	// Resolve the UDP address for the specific port on which we expect broadcasts
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("error resolving UDP address: %w", err)
	}

	// Create a UDP connection for listening
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("error listening to UDP port %d: %w", port, err)
	}
	defer conn.Close()

	// Buffer to store incoming messages
	buf := make([]byte, 1024)

	// Listen for incoming broadcast messages
	fmt.Println("Listening for broadcast messages...")
	for {
		n, senderAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return fmt.Errorf("error reading from UDP: %w", err)
		}

		// Print the received broadcast message
		fmt.Printf("Received broadcast from %s: %s\n", senderAddr, string(buf[:n]))
	}
}

func main() {
	err := listenForBroadcast(6000)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
