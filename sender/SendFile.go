package sender

import (
	"fmt"
	"net"
)

// SendFile takes in a fileobject and sends the file over
// a given connection.
func SendFile(filepath string, conn net.TCPConn) error {
	var encodedFile = prepFile(filepath)
	writeLength, err := conn.Write(encodedFile)
	if err != nil {
		fmt.Printf("unable to send file: \n", err)
		return err
	} else {
		fmt.Printf("File written successfully: \n %d bytes in total", writeLength)
	}
	return nil
}
