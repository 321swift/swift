package main

import (
	"fmt"
	"os"
	"swift/sender"
)

// var wg sync.WaitGroup

func main() {
	sender.Send()
	fmt.Println(os.Hostname())
}
