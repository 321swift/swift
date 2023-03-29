package utils

import (
	"fmt"
	"net"
)

// SendMessage sends a message on a given IP address and port number
func SendMessage(address string, port int, message string) error {
	// convert message to a byte array
	messageInBytes := []byte(message)

	// Resolve the IP address and port
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return err
	}

	// Create the UDP socket
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send the message
	_, err = conn.Write(messageInBytes)
	if err != nil {
		return err
	}

	return nil
}
