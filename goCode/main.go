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

	go sender()
	message, senderAddress, _ := c.ListenForBroadcastMessage(5050)
	fmt.Println(string(message), senderAddress)

	// fmt.Println(u.GetIp())
}

func sender() {
	defer wg.Done()

	// // send message on the broadcast address with port 5050
	u.GetIp()
	for i := 0; i < 10; i++ {
		c.SendMessage(u.GetBroadcastAddress(), 5050, "hello i am swift")
		fmt.Printf("sent %d \n", i)
		time.Sleep(time.Second * 3)
	}

}
