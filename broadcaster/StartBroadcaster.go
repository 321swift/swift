package broadcaster

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func StartupBroadcaster() {
	var wg sync.WaitGroup

	// get available port number
	availablePorts, err := getAvailablePorts()
	if err != nil {
		fmt.Println("Error getting available ports:", err.Error())
		return
	}

	var socketPort int
	var broadcastPort int

	randNum := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		socketPort = availablePorts[randNum.Intn(4000)]
		broadcastPort = availablePorts[randNum.Intn(4000)]

		if socketPort != broadcastPort {
			break
		}
	}

	wg.Add(2)
	go startSocket(socketPort)
	go sendBroadcast(broadcastPort, socketPort)

	wg.Wait()

}
