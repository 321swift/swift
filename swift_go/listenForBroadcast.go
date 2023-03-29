package main

import (
	"fmt"
	"net"
	"time"
)

// ListenForBroadcast listens on the network for a broadcasted message
func ListenForBroadcast(port int) ([]byte, error) {
	// Resolve the broadcast address and port
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	// Create the UDP socket
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Set a timeout for the socket
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	// Wait for a message
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}
