package main

import (
	"fmt"
	"net"
	c "swift/core"
	u "swift/utils"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	// wg.Add(1)

	// go sender()
	// message, senderAddress, _ := c.ListenForBroadcastMessage(5050)
	// fmt.Println(string(message), senderAddress)
	addrs, _ := net.InterfaceAddrs()
	print(addrs[0].String())
	for _, ad := range addrs {
		fmt.Println(ad.Network(), ad.String())
	}
	str, _ := u.CalcBroadcastAddress("10.20.84.239/21")
	print(str)

}

func sender() {
	defer wg.Done()

	// send message on the broadcast address with port 5050
	u.GetIp()
	for i := 0; i < 10; i++ {
		c.SendMessage(u.GetBroadcastAddress(), 5050, "Hello, I am swift")
		fmt.Printf("sent %d \n", i)
		time.Sleep(time.Second * 3)
	}

}
