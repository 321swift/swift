package broadcaster

import (
	"log"
	"net"
)

func getAllBroadcasts(ips []net.Addr) []string {
	broadcasts := make([]string, 0)

	for _, ip := range ips {
		broadcast, err := calcBroadcastAddress(ip.String())
		if err != nil {
			log.Fatal("error while calculating broadcast Address from : ", ip)
		}
		broadcasts = append(broadcasts, broadcast)
	}
	return broadcasts
}
