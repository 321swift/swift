package receiver

import (
	"net"
	"sync"
)

var wg sync.WaitGroup

func ListenRoutine(port int) ([]byte, net.Addr, error) {
	wg.Done()

	return Listen(port)
}
