package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"swift2/global"
	"time"
)

type server struct {
	conn          net.Conn
	hostname      string
	serverPort    int
	listener      net.Listener
	serverTimeout time.Duration
	timer         *time.Timer
	logger        *global.Logger
}

func NewServer(logger *global.Logger) *server {
	name, _ := os.Hostname()
	sPort := global.GetAvailablePort()
	global.BackendServerPort = sPort
	return &server{
		hostname:      name,
		serverPort:    sPort,
		serverTimeout: time.Second * 10,
		logger:        logger,
	}
}

func (s *server) Shutdown() {
	defer s.logger.WriteLog("Server stopped")
	s.listener.Close()
	if s.conn != nil {
		err := s.conn.Close()
		if err != nil {
			return
		} else {
			s.logger.WriteLog(fmt.Sprint("unable to close connection: ", err))
			return
		}
	}
}

func (s *server) Start() {
	go func() {
		s.Broadcast()
	}()

	s.timer = time.AfterFunc(s.serverTimeout, func() {
		s.logger.WriteLog("Server timeout... shutting down")
		s.Shutdown()
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.serverPort))
	if err != nil {
		log.Fatal(err)
	}
	s.listener = listener

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.WriteLog(fmt.Sprint("accepting connection err: ", err))
			break
		}
		s.timer.Stop()
		s.logger.WriteLog(fmt.Sprint("Connection made: ", conn))
		s.conn = conn
		conn.Write([]byte("Welcome to this "))
		go s.readLoop(conn)
	}
}

func (s *server) readLoop(conn net.Conn) {
	defer conn.Close()

	for {
		// receive file size
		var dataSize int64
		err := binary.Read(conn, binary.LittleEndian, &dataSize)
		if err != nil {
			log.Fatal(err)
		}

		// recieve data prefixed with filename
		data := new(bytes.Buffer)
		i, err := io.CopyN(data, conn, dataSize)
		if err != nil {
			log.Fatal(err)
		}
		s.logger.WriteLog(fmt.Sprintf("Received %d bytes from %v", i, conn.RemoteAddr()))

		// seperate data from filename
		filename, fileContent, ok := bytes.Cut(data.Bytes(), []byte("$$$$"))
		if !ok {
			s.logger.WriteLog("Unable to parse file... ")
		}

		err = os.WriteFile(fmt.Sprintf("./1%s", filename), fileContent, os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}
	}

}

func (s *server) Send(filePath string) error {
	fmt.Println("sending")
	data, err := os.ReadFile(filePath)
	if err != nil {
		s.logger.WriteLog(err)
		return err
	}

	// prefix data with filename
	filename := []byte(fmt.Sprintf("%v$$$$", path.Base(filePath)))
	data = append(filename, data...)

	// send file size
	go func() {
		err = binary.Write(s.conn, binary.LittleEndian, int64(len(data)))
		if err != nil {
			s.logger.WriteLog(err)
			return
		}

		// send data
		i, err := io.CopyN(s.conn, bytes.NewReader(data), int64(len(data)))
		if err != nil {
			s.logger.WriteLog(err)
			return
		}
		s.logger.WriteLog(fmt.Sprintf("File sent successfully: %d / %d bytes written", i, len(data)))
	}()

	return nil
}

// func (s *server) Receive() {

// 	// receive file size
// 	var dataSize int64
// 	err := binary.Read(s.conn, binary.LittleEndian, &dataSize)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// recieve data prefixed with filename
// 	data := new(bytes.Buffer)
// 	i, err := io.CopyN(data, s.conn, dataSize)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("Received %d bytes from connection", i)

// 	// seperate data from filename
// 	filename, fileContent, ok := bytes.Cut(data.Bytes(), []byte("$$$$"))
// 	if !ok {
// 		log.Println("Unable to parse file... ")
// 	}

// 	err = os.WriteFile(fmt.Sprintf("./1%s", filename), fileContent, os.ModePerm)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
