package node

import "net"

type Node struct {
	connectionPool []int
	hostname       string
	uiConnection   net.Conn
}

// funcs

func (n *Node) getAvailablePort() int {
	return 0
}

func (n *Node) setupConnectionPool() {

}

func (n *Node) logInfo() {}
func (n *Node) logErr()  {}
