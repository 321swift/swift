package utils

import "fmt"

func GetBroadcastAddress() string {
	fmt.Println("Begin obtaining broadcast address")
	interfaces, err := GetUpnRunninginterfaces()
	if err != nil {
		panic("error while getting up and runnign interfaces")
	}
	fmt.Println("Extracting ipv4 Address")
	addr, _ := ExtractIPV4Address(interfaces[0])

	fmt.Println("Calculating broadcast")
	addr, _ = CalcBroadcastAddress(addr)
	return addr
}
