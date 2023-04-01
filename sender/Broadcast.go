package sender

import (
	"fmt"
	"net"
	"os"
	"swift/core"
	"swift/utils"
	"time"
)

func Broadcast(port int, serverPort int) {
	ips, _ := net.InterfaceAddrs()
	ips = utils.FilterIPs(ips)

	broadcasts := utils.GetAllBroadcasts(ips)
	devHostname, err := os.Hostname()
	if err != nil {
		devHostname = "user"
	}

	for i := 15; i > 0; i-- {
		for _, addr := range broadcasts {
			core.SendMessage(
				addr,
				port,
				fmt.Sprintf("%s@swift:%d", devHostname, serverPort),
			)
		}
		time.Sleep(time.Second * 1)
		fmt.Println("Broadcast sent ", i)
	}
}
