package utils

func GetBroadcastAddress() string {
	interfaces, err := GetUpnRunninginterfaces()
	if err != nil {
		panic("error while getting up and runnign interfaces")
	}
	addr, _ := ExtractIPV4Address(interfaces[0])

	addr, _ = CalcBroadcastAddress(addr)
	return addr
}
