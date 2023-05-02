package node

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
)

func (n *Node) GetAvailablePort() int {
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
