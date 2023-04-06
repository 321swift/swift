package main

import (
	"fmt"
	"swift/listener"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	var ch = make(chan string, 1024)
	go func() {
		defer wg.Done()
		err := listener.Listener(ch)
		if err != nil {
			fmt.Println("unable to listen: ", err)
		}
	}()
	endTime := time.Now().Add(time.Second * 10)

	for time.Now().Before(endTime) {
		msg := <-ch
		if msg != "" {
			fmt.Println(msg)
		}
	}
	wg.Wait()
}
