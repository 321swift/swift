package main

import (
	"fmt"
	"swift2/backend/client"
	"swift2/global"
	"swift2/ui"
	"time"
)

func main() {
	// create new listener
	client := client.NewClient()
	time.AfterFunc(time.Second*8, func() {
		client.Connect(fmt.Sprintf(":%d", global.BackendServerPort))
	})

	// start server
	server := ui.NewUiServer()
	server.Start()
}
