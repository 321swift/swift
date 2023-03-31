package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"swift/sender"
	"sync"
)

var wg sync.WaitGroup

func main() {
	role := WelcomeScreen()
	if role == "1" {
		//call the sender function
	} else if role == "2" {
		// call the receiver function
	}
}

func WelcomeScreen() string {
	var str string

	for i := false; !i; {
		fmt.Println("")
		fmt.Println("Welcome to the swift application")
		fmt.Println("Select the role you wish to assume: \n\t 1. Sending \n\t 2. Receiving")

		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println(err)
		} else if len(line) != 1 &&
			!strings.ContainsAny(line, "12") {
			fmt.Println("Please enter one number only acording to the given options")
		} else {
			str = line
			i = true
		}
	}
	return str
}

func BroadcastRoutine(port int, serverPort int) {
	defer wg.Done()
	sender.Broadcast(port, serverPort)
}

func ServerRouting(port int) {
	defer wg.Done()
	sender.StartServer(port)
}

func ConnectToSocket(ip string, port string) (net.Conn, error) {
	// Connect to the server using the specified IP and port
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	println("Connection made successfully")
	return conn, nil
}
