package broadcaster

import (
	"net"
	"strings"
)

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
