package broadcaster

import (
	"fmt"
	"net"
	"os"
	"time"
)

func SendBroadcast(broadcastPort int, socketPort int) {
	ips, _ := net.InterfaceAddrs()
	ips = filterIPs(ips)

	broadcasts := getAllBroadcasts(ips)
	devHostname, err := os.Hostname()
	if err != nil {
		devHostname = "user"
	}

	for i := 15; i > 0; i-- {
		for _, addr := range broadcasts {
			sendMessage(
				addr,
				broadcastPort,
				fmt.Sprintf("%s@swift:%d", devHostname, socketPort),
			)
		}
		time.Sleep(time.Second * 1)
		fmt.Println("Broadcast sent ", i)
	}
}
