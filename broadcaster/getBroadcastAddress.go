package broadcaster

import "fmt"

func getBroadcastAddress() string {
	interfaces, err := getUpnRunninginterfaces()
	if err != nil {
		panic("error while getting up and runnign interfaces")
	}
	fmt.Println(interfaces)
	addr, _ := extractIPV4Address(interfaces[0])

	addr, _ = calcBroadcastAddress(addr)
	return addr
}
