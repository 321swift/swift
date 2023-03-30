package sender

import "sync"

var wg sync.WaitGroup

func SendRoutine() {
	defer wg.Done()
	Send()
}
