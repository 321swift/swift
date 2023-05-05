package node

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	BroadcastPort = 51413
)

func (n *Node) getAvailablePort() int {
	desiredPort := n.basePort
	for {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", desiredPort))
		if err == nil {
			listener.Close()
			return desiredPort
		}
		// Port is already in use, so try the next one
		desiredPort++

	}
}

func OpenPage(url string) {
	var err error

	switch runtime.GOOS {
	case "darwin":
		err = exec.Command("open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}

	if err != nil {
		fmt.Println("Failed to open URL:", err)
	}
}

func (n *Node) handleFileReception(w http.ResponseWriter, r *http.Request) {
	filename := r.Header.Get("filename")
	chunkSize := r.Header.Get("chunkSize")
	chunks := r.Header.Get("chunks")
	n.infoLog.Printf(filename, chunkSize, chunks)

	// upgrade connection to websocket
	// var upgrader = websocket.Upgrader{
	// 	ReadBufferSize:  1024,
	// 	WriteBufferSize: 1024,
	// }

	// upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// conn, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// messageType, message, err := conn.ReadMessage()

}

func (n *Node) broadcast() {
	ipnets := getActiveIps()
	addrs := make([]net.IP, 0)
	for _, inet := range ipnets {
		ip, err := calcBroadcastAddress(inet)
		if err != nil {
			continue
		}
		addrs = append(addrs, ip)
	}
	n.infoLog.Printf("Sending broadcast to %v\n", addrs)

	for i := 0; i < 10; i++ {
		for _, addr := range addrs {
			go func(addr net.IP) {
				sendMessage(addr, fmt.Sprintf("%s:%d", "swift", n.serverPort))
			}(addr)
		}
		time.Sleep(time.Second)
	}
	n.infoLog.Println("log", "Broadcast sent out")
}

func getActiveIps() []net.IPNet {
	interfaces, err := getUpnRunninginterfaces()
	var addrs []net.IPNet
	if err != nil {
		panic("error while getting up and runnign interfaces")
	}
	for _, interf := range interfaces {
		addr := extractIPV4Address(interf)
		if addr != nil {
			addrs = append(addrs, *addr)
		}
	}
	return addrs
}

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

func calcBroadcastAddress(ipSub net.IPNet) (net.IP, error) {
	// Parse the IP address and subnet
	ip, ipNet, err := net.ParseCIDR(ipSub.String())
	if err != nil {
		return nil, err
	}

	// Get the network size in bits
	ones, bits := ipNet.Mask.Size()

	// Calculate the broadcast address
	mask := net.CIDRMask(ones, bits)
	network := ip.Mask(mask)
	broadcast := make(net.IP, len(network))
	for i := range network {
		broadcast[i] = network[i] | ^mask[i]
	}

	return broadcast, nil
}

func sendMessage(address net.IP, message string) error {
	// convert message to a byte array
	messageInBytes := []byte(message)

	// Resolve the IP address and port
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address.String(), BroadcastPort))
	if err != nil {
		return err
	}

	// Create the UDP socket
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send the message
	_, err = conn.Write(messageInBytes)
	if err != nil {
		return err
	}

	return nil
}

func extractIPV4Address(iface net.Interface) *net.IPNet {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil
	}
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet
		}
	}
	return nil
}

func (n *Node) sendIntroduction() error {
	keys := make([]int, 0)
	for i := range n.connectionPool {
		keys = append(keys, i)
	}
	response := Intro{
		Hostname: n.hostname,
		Conns:    keys,
	}
	err := json.NewEncoder(n.backendConnection).Encode(&response)
	if err != nil {
		return err
	}
	return nil
}
func (n *Node) receiveIntroduction() (Intro, error) {
	// Receive the introduction message
	var intro Intro
	err := json.NewDecoder(n.backendConnection).Decode(&intro)
	if err != nil {
		return Intro{}, err
	}
	intro.Status = "connected"
	intro.ConnectedNodeIP = strings.Split(n.backendConnection.RemoteAddr().String(), ":")[0]

	return intro, nil
}
