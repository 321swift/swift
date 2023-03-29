package utils

import "net"

// GetBroadcastAddress calculates the broadcast address for a given IP address and subnet
func GetBroadcastAddress(ipAddress string) (string, error) {
	// Parse the IP address and subnet
	ip, ipNet, err := net.ParseCIDR(ipAddress)
	if err != nil {
		return "", err
	}

	// Get the network size in bits
	ones, bits := ipNet.Mask.Size()

	// Calculate the broadcast address
	mask := net.CIDRMask(ones, bits)
	network := ip.Mask(mask)
	broadcast := make(net.IP, len(network))
	for i := range network {
		broadcast[i] = network[i] | ^mask[i]
	}

	return broadcast.String(), nil
}
