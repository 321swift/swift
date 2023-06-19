package global

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
)

const (
	BroadcastPort = 51413
)

var BackendServerPort = 0

func GetAvailablePort(desiredPort int) int {
	for {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", desiredPort))
		if err == nil {
			listener.Close()
			return desiredPort
		}
		// Port is already in use, so try the next one
		log.Printf("Port %d already in use, trying next port\n", desiredPort)
		desiredPort++

	}
}

func CreateDirectoryIfNotExists(dirName string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "/"
	}

	dirPath := filepath.Join(homeDir, dirName)

	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0755)
		if err != nil {
			return "/"
		}
	}

	return dirPath
}

type Message struct {
	Filename string `json:"Filename"`
	Data     []byte `json:"Data"`
}
