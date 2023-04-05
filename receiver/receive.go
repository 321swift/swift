package receiver

import "net"

func receiveFile(conn net.TCPConn) ([]byte, error) {
	var fileStream []byte
	_, err := conn.Read(fileStream)
	if err != nil {
		return nil, err
	}

	return fileStream, nil
}
