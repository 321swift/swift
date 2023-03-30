package sender

import (
	"fmt"
	"net"
	"os"
	"swift/core"
	"swift/utils"
	"time"
)

func Send() {
	ips, _ := net.InterfaceAddrs()
	ips = utils.FilterIPs(ips)

	broadcasts := utils.GetAllBroadcasts(ips)
	devHostname, err := os.Hostname()
	if err != nil {
		devHostname = "user"
	}

	for i := 15; i > 0; i-- {
		for _, addr := range broadcasts {
			core.SendMessage(addr, 5050, fmt.Sprintf("%s@swift", devHostname))
		}
		time.Sleep(time.Second * 1)
		fmt.Println("Sender sent ", i)
	}
}
