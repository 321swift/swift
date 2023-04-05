package broadcaster

import (
	"net"
	"strings"
)

func filterIPs(ips []net.Addr) []net.Addr {
	filtered := make([]net.Addr, 0)
	for _, ip := range ips {
		if !strings.ContainsAny(ip.String(), "abcdef:") &&
			!net.ParseIP(strings.Split(ip.String(), "/")[0]).IsLoopback() {
			filtered = append(filtered, ip)
		}
	}

	return filtered
}
