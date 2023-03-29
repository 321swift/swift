package main

import (
	"fmt"
	"net"
	"strings"
)

// GetLocalIPAddress returns the IP address of the current network
func GetLocalIPAddress() (net.IP, net.IPMask, error) {
	// Get the default route
	route, err := net.InterfaceAddrs()
	if err != nil {
		return nil, nil, err
	}

	// Loop through the interfaces and find the default route
	for _, r := range route {
		if ipnet, ok := r.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, ipnet.Mask, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("no default route found")
}

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

func getUpnRunninginterfaces() ([]net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var upnRunning []net.Interface

	// get all up and running interfaces
	for _, interf := range interfaces {
		if flags := interf.Flags; strings.Contains(flags.String(), "running") &&
			strings.Contains(flags.String(), "up") &&
			!strings.Contains(interf.Name, "VirtualBox") &&
			!strings.Contains(interf.Name, "Loopback") {
			upnRunning = append(upnRunning, interf)
		}
	}
	return upnRunning, nil
}

func ExtractIPV4Address(iface net.Interface) (string, error) {

	address, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	return address[1].String(), nil
}

func main() {
	interfaces, err := getUpnRunninginterfaces()
	if err != nil {
		panic("error while getting up and runnign interfaces")
	}
	addr, _ := ExtractIPV4Address(interfaces[0])
	fmt.Println(addr)
}
