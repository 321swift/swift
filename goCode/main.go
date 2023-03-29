package main

import (
	"fmt"
	u "swift/utils"
)

func main() {
	interfaces, err := u.GetUpnRunninginterfaces()
	if err != nil {
		panic("error while getting up and runnign interfaces")
	}
	addr, _ := u.ExtractIPV4Address(interfaces[0])
	fmt.Println(addr)
}
