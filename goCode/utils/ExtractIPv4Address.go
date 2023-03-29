package utils

import "net"

func ExtractIPV4Address(iface net.Interface) (string, error) {

	address, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	return address[1].String(), nil
}
