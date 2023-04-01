package main

import (
	"swift/receiver"
	"swift/sender"
	"swift/utils"
	"sync"
)

func main() {
	role := utils.WelcomeScreen()
	var wg sync.WaitGroup

	serverPort := 1234
	if role == 1 {
		sender.StartSender(&wg, serverPort)
	} else if role == 2 {
		receiver.StartReceiver(&wg)
	}
	wg.Wait()
}
