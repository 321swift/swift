package broadcaster

import (
	"fmt"
	"net"
	"os"
	"time"
)

func sendBroadcast(broadcastPort int, socketPort int) {
	ips, _ := net.InterfaceAddrs()
	ips = filterIPs(ips)

	broadcasts := getAllBroadcasts(ips)
	devHostname, err := os.Hostname()
	if err != nil {
		devHostname = "user"
	}

	message := fmt.Sprintf("%s@swift:%d", devHostname, socketPort)
	fmt.Printf("Sending Broadcast: %v \n", message)

	for i := 15; i > 0; i-- {
		for _, addr := range broadcasts {
			sendMessage(
				addr,
				broadcastPort,
				message,
			)
		}
		time.Sleep(time.Second * 1)
	}
}
