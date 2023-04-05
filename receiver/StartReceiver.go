package receiver

import (
	"fmt"
	"net"
)

func StartReceiver(conn net.TCPConn) error {
	file, err := receiveFile(conn)
	if err != nil {
		fmt.Println(err)
		return err
	}

	decryptedFile, err := decryptFile(file)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = saveFile("received", "./", decryptedFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
