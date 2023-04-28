package server

import (
	"fmt"
	"net"
	"strings"
	"swift2/global"
	"sync"
	"time"
)

func (s *Server) Broadcast() {
	var wg sync.WaitGroup
	ipnets := s.GetActiveIps()
	addrs := make([]net.IP, 0)
	for _, inet := range ipnets {
		ip, err := s.calcBroadcastAddress(inet)
		if err != nil {
			continue
		}
		addrs = append(addrs, ip)
	}
	s.logger.Write("log", fmt.Sprintf("Sending broadcast to %v", addrs))

	wg.Add(20 * len(addrs))
	for i := 0; i < 10; i++ {
		for j, addr := range addrs {
			go func(addr net.IP, offset int) {
				defer wg.Done()
				s.sendMessage(addr, fmt.Sprintf("%s:%d", s.hostname, s.serverPort))
			}(addr, j+i)
		}
		time.Sleep(time.Millisecond * 100)
	}
	s.logger.Write("log", "Broadcast sent out")
	wg.Wait()
}

func (s *Server) GetActiveIps() []net.IPNet {
	interfaces, err := s.getUpnRunninginterfaces()
	var addrs []net.IPNet
	if err != nil {
		panic("error while getting up and runnign interfaces")
	}
	for _, interf := range interfaces {
		addr := s.extractIPV4Address(interf)
		if addr != nil {
			addrs = append(addrs, *addr)
		}
	}
	return addrs
}

// GetBroadcastAddress calculates the broadcast
// address for a given IP address and subnet
// in the format: ip/subnet
func (s *Server) calcBroadcastAddress(ipSub net.IPNet) (net.IP, error) {
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

func (s *Server) extractIPV4Address(iface net.Interface) *net.IPNet {
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

func (s *Server) getUpnRunninginterfaces() ([]net.Interface, error) {
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

// SendMessage sends a message on a given IP address and port number
func (s *Server) sendMessage(address net.IP, message string) error {
	// convert message to a byte array
	messageInBytes := []byte(message)

	// Resolve the IP address and port
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address.String(), global.BroadcastPort))
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
