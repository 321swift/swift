package broadcaster

import (
	"net"
)

func extractIPV4Address(iface net.Interface) (string, error) {

	address, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	// fmt.Println("address\n", address)
	return address[len(address)-1].String(), nil
}
