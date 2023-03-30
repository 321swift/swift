package sender

import (
	"fmt"
	"net"
	"os"
	"swift/core"
	"swift/utils"
)

func Send() {
	ips, _ := net.InterfaceAddrs()
	ips = utils.FilterIPs(ips)

	broadcasts := utils.GetAllBroadcasts(ips)

	devHostname, err := os.Hostname()

	for _, addr := range broadcasts {
		if err != nil {

			core.SendMessage(addr, 5050, "@swift swiftUser")
		}
		core.SendMessage(addr, 5050, fmt.Sprintf("@swift %s", devHostname))
	}
}
