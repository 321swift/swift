package core

import (
	"fmt"
	"net"
	"time"
)

// // ListenForBroadcastMessage listens on the given broadcast address and port for a message
// func ListenForBroadcast(broadcastAddress string, port int) ([]byte, error) {
// 	// Resolve the broadcast address and port
// 	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", broadcastAddress, port))
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create a UDP socket to listen on
// 	conn, err := net.ListenUDP("udp", addr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer conn.Close()

// 	// Set a timeout for the socket
// 	conn.SetReadDeadline(time.Now().Add(time.Second * 15))

// 	// Wait for a message
// 	buffer := make([]byte, 1024)
// 	fmt.Print("Listening for message \n ")
// 	n, _, err := conn.ReadFromUDP(buffer)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return buffer[:n], nil
// }

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
