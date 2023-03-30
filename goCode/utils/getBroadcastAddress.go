package utils

import "fmt"

func GetBroadcastAddress() string {
	interfaces, err := GetUpnRunninginterfaces()
	if err != nil {
		panic("error while getting up and runnign interfaces")
	}
	fmt.Println(interfaces)
	addr, _ := ExtractIPV4Address(interfaces[0])

	addr, _ = CalcBroadcastAddress(addr)
	return addr
}
