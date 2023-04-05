package listener

import (
	"fmt"
	"net"
	"time"
)

func StartListener(broadcastPort int64) (string, *net.UDPAddr, error) {
	// Resolve the broadcast address and port
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", broadcastPort))

	if err != nil {
		return "", nil, err
	}

	// Create a UDP socket to listen on
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return "", nil, err
	}
	defer conn.Close()

	// Set a timeout for the socket
	conn.SetReadDeadline(time.Now().Add(time.Second * 50))

	println("now listening on port ", broadcastPort)
	// Wait for a message
	var buffers string

	var remoteAddr *net.UDPAddr
	var n int

	for {

		buffer := make([]byte, 1024)
		n, remoteAddr, err = conn.ReadFromUDP(buffer)
		if err != nil {
			return "", nil, err
		}

		buffers = fmt.Sprintf("%s+%s", buffers, buffer[:n])
		if conn == nil {
			break
		}
	}

	return buffers, remoteAddr, nil
}
