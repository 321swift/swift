package main

import (
	"fmt"
	c "swift/core"
	u "swift/utils"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)

	// // listen for message on the broadcast address
	// message := make(chan []byte)
	// go listener(message)
	// fmt.Print(string(<-message))

	// // send message on the broadcast address with port 5050
	// for i := 0; i < 10; i++ {
	// 	c.SendMessage(u.GetBroadcastAddress(), 5050, "hi there I am David")
	// 	fmt.Printf("sent %d \n", i)
	// 	time.Sleep(time.Second * 3)
	// }
	// close(message)
	// c.ListenForBroadcast()

	go sender()
	message, senderAddress, _ := c.ListenForBroadcastMessage(5050)
	fmt.Println(string(message), senderAddress)

}

// func listener(message chan []byte) {
// 	defer wg.Done()

// 	bytes, _ := c.ListenForBroadcast(5050)
// 	message <- bytes
// }

func sender() {
	defer wg.Done()

	// // send message on the broadcast address with port 5050
	u.GetIp()
	for i := 0; i < 10; i++ {
		c.SendMessage(u.GetBroadcastAddress(), 5050, "hi there I am David")
		fmt.Printf("sent %d \n", i)
		time.Sleep(time.Second * 3)
	}

}
