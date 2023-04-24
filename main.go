package main

import (
	"swift2/ui"
	"time"

	"golang.org/x/vuln/client"
)

func main() {
	// create new listener
	client := client.NewClient()
	time.AfterFunc(time.Second*8, func() {
client.
	})
	
	// start server
	server := ui.NewUiServer()
	server.Start()
}
