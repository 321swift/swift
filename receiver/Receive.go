package receiver

import "net"

func ReceiveFile(conn net.TCPConn) ([]byte, error) {
	var fileStream []byte
	_, err := conn.Read(fileStream)
	if err != nil {
		return nil, err
	}

	return fileStream, nil
}
