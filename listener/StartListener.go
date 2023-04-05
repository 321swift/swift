package listener

import (
	"net"
)

// // the start listener method listens for a broadcast
// // on a given port and sends all the messages it receives
// // to the channel it is given.
// func StartListener(broadcastPort int64, channel chan ListenerChannel) {
// 	defer close(channel)

// 	// Resolve the broadcast address and port
// 	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", broadcastPort))
// 	if err != nil {
// 		log.Fatal("unable to resolve address: ", err)
// 	}

// 	// Create a UDP socket to listen on
// 	conn, err := net.ListenUDP("udp", addr)
// 	if err != nil {
// 		log.Fatal("Unable to setup listener", err)
// 	}
// 	defer conn.Close()

// 	// Set a timeout for the socket
// 	conn.SetReadDeadline(time.Now().Add(time.Second * 2))

// 	println("now listening on port ", broadcastPort)
// 	// Wait for a message
// 	var buffers string

// 	var remoteAddr *net.UDPAddr
// 	var n int

// 	buffer := make([]byte, 1024)
// 	n, remoteAddr, err = conn.ReadFromUDP(buffer)

// 	if err != nil {
// 		log.Fatal("error while listening for broadcast", err)
// 	}

// 	buffers = fmt.Sprintf("%s+%s", buffers, buffer[:n])
// 	channel <- ListenerChannel{buffers, remoteAddr}
// }

// listenAndServe listens for a broadcast on the given port
// and sends all the messages it receives to the given channel.
func ListenAndServe(port string, messageChan chan<- string) error {
	// Resolve the UDP address.
	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		return err
	}

	// Listen for UDP packets on the port.
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create a buffer for receiving packets.
	buf := make([]byte, 4096)

	// Continuously read packets and send their contents to the channel.
	for {
		n, senderAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		// Extract the message from the packet and send it to the channel.
		message := string(buf[:n]) + "-" + senderAddr.String()
		messageChan <- message
	}
}
