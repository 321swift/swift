package sender

import "sync"

func StartSender(wg *sync.WaitGroup, serverPort int) {
	wg.Add(2)

	// start broadcast
	go func() {
		defer wg.Done()
		Broadcast(5050, serverPort)
	}()

	// start server
	go func() {
		defer wg.Done()
		StartServer(serverPort)
	}()
}
