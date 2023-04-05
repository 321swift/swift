package main

import (
	"fmt"
	"swift/listener"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	var channel = make(chan string)

	go func() {
		listener.ListenAndServe("5050", channel)
		wg.Done()
	}()

	go func() {
		var msg string
		for {
			msg = <-channel
			if msg != "" {
				fmt.Println(msg)
			} else {
				break
			}
		}
	}()

	wg.Wait()
}
