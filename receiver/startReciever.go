package receiver

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

func StartReceiver(wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		// start listener

		defer wg.Done()
		bt, addr, err := Listen(5050)
		if err != nil {
			log.Fatalf("could not listen: %v", err)
		}
		fmt.Println(string(bt))
		fmt.Println(addr.String())

		connectAddr := fmt.Sprintf("%v:%v",
			strings.Split(addr.String(), ":")[0], //obtain ip
			strings.Split(string(bt), ":")[1],    // obtain port
		)

		ConnectToSocket(connectAddr)

	}()
}
