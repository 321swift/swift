package main

import (
	"fmt"
	"net"
	"strings"
)

// GetBroadcastAddress returns the broadcast address of the current network
func GetBroadcastAddress() (net.IP, error) {
	// Get the default route
	route, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	// Loop through the interfaces and find the default route
	for _, r := range route {
		if ipnet, ok := r.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// Calculate the broadcast address using the subnet mask
				broadcastIP := make(net.IP, len(ipnet.IP))
				copy(broadcastIP, ipnet.IP)
				for i := range broadcastIP {
					broadcastIP[i] |= ^ipnet.Mask[i]
				}
				return broadcastIP, nil
			}
		}
	}
	return nil, fmt.Errorf("no default route found")
}

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

func PrintAllIp() {
	route, err := net.InterfaceAddrs()
	fmt.Print(route, err)
}

func main() {
	print := fmt.Println
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("an error occurred")
	}
	// get all up and running interfaces
	for _, interf := range interfaces {
		if flags := interf.Flags; strings.Contains(flags.String(), "running") &&
			strings.Contains(flags.String(), "up") &&
			!strings.Contains(interf.Name, "VirtualBox") &&
			!strings.Contains(interf.Name, "Loopback") {
			print(interf)
		}
	}
}
