package receiver

import (
	"fmt"
	"net"
	"time"
)

func Listen(port int) ([]byte, net.Addr, error) {
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
	conn.SetReadDeadline(time.Now().Add(time.Second * 35))

	println("now listening on port ", port)
	// Wait for a message
	buffer := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, nil, err
	}

	return buffer[:n], remoteAddr, nil
}
