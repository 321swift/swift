package broadcaster

import (
	"fmt"
	"net"
)

func getAvailablePorts() ([]int, error) {
	var availablePorts []int

	// Get list of network interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("error getting network interfaces: %s", err)
	}

	// Loop over network interfaces
	for _, iface := range ifaces {
		// Get list of unicast addresses for interface
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("error getting addresses for interface %s: %s", iface.Name, err)
		}

		// Loop over unicast addresses
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				// Check if address is IPv4 and not a loopback address
				if v.IP.To4() != nil && !v.IP.IsLoopback() {
					// Scan for free ports on address
					for port := 3000; port <= 8000; port++ {
						conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", v.IP.String(), port), 10)
						if err != nil {
							availablePorts = append(availablePorts, port)
						} else {
							conn.Close()
						}
					}
				}
			}
		}
	}

	return availablePorts, nil
}
