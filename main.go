package main

import (
	"fmt"
	"swift/receiver"
	"swift/sender"
	"sync"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)

	portNumber := 5050
	go BroadcastRoutine(portNumber)
	message, addr, err := receiver.Listen(portNumber)

	fmt.Println(string(message))
	fmt.Println(addr)
	fmt.Println(err)
}

func BroadcastRoutine(port int) {
	defer wg.Done()
	sender.Broadcast(port)
}
