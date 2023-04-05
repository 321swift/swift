package broadcaster

func GetIp() []string {

	interfaces, err := getUpnRunninginterfaces()
	var addrs []string
	if err != nil {
		panic("error while getting up and runnign interfaces")
	}
	for _, interf := range interfaces {
		addr, _ := extractIPV4Address(interf)
		addrs = append(addrs, addr)
	}
	return addrs
}
