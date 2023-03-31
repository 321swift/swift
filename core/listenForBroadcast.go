package core

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
	conn.SetReadDeadline(time.Now().Add(time.Second * 15))

	// Wait for a message
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:n], nil
}

// ListenForBroadcastMessage listens on the given broadcast address and port for a message
// and returns the message and the IP address of the sender
func ListenForBroadcastMessage(port int) ([]byte, net.Addr, error) {
	// Resolve the broadcast address and port
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, nil, err
	}

	// Create a UDP socket to listen on
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	// Set a timeout for the socket
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	// Wait for a message
	buffer := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, nil, err
	}

	return buffer[:n], remoteAddr, nil
}
