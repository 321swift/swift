package main

import (
	"fmt"
	"swift/receiver"
	"swift/sender"
	"sync"
)

var wg sync.WaitGroup

func main() {
	role := WelcomeScreen()
	// fmt.Println(role, reflect.TypeOf(role))
	if role == 1 {
		// 	//call the sender function
		sender.StartServer(5053)

		// fmt.Println("role == 1")
	} else if role == 2 {
		// 	// call the receiver function
		receiver.Listen(5050)
		// fmt.Println("role == 2")
	}
}

func WelcomeScreen() int64 {
	var role int64

	for i := false; !i; {
		fmt.Println("")
		fmt.Println("Welcome to the swift application")
		fmt.Println("Select the role you wish to assume: \n\t 1. Sending \n\t 2. Receiving")

		var number int64
		_, err := fmt.Scanf("%d", &number)

		if err != nil {
			fmt.Println(err)
		} else if number < 0 {
			fmt.Println("Please enter one number only acording to the given options")
		} else if number > 2 {
			fmt.Println("Please enter one number only acording to the given options")
		} else {
			i = true
			role = number
		}
	}
	return role
}

func ServerRoutine(port int, broadcastPort int) {
	wg.Add(1)
	go startBroadcast(broadcastPort, port)
	sender.StartServer(port)
}

func startBroadcast(broadcastPort int, port int) {
	defer wg.Done()
	sender.Broadcast(broadcastPort, port)
}

// func ConnectToSocket(ip string, port string) (net.Conn, error) {
// 	// Connect to the server using the specified IP and port
// 	conn, err := net.Dial("tcp", ip+":"+port)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}

// 	println("Connection made successfully")
// 	return conn, nil
// }
