package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"swift2/global"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	client        net.Conn
	hostname      string
	serverPort    int
	listener      net.Listener
	serverTimeout time.Duration
	timer         *time.Timer
	logger        *global.Logger
}

func NewServer(logger *global.Logger) *Server {
	name, _ := os.Hostname()
	sPort := global.GetAvailablePort(5150)
	global.BackendServerPort = sPort
	return &Server{
		hostname:      name,
		serverPort:    sPort,
		serverTimeout: time.Second * 10,
		logger:        logger,
	}
}

func (s *Server) Shutdown() {
	defer s.logger.WriteLog("Server stopped")
	s.listener.Close()
	if s.client != nil {
		err := s.client.Close()
		if err != nil {
			return
		} else {
			s.logger.WriteLog(fmt.Sprint("unable to close connection: ", err))
			return
		}
	}
}

func (s *Server) Start() {
	go func() {
		s.Broadcast()
	}()

	s.timer = time.AfterFunc(s.serverTimeout, func() {
		s.logger.WriteLog("Server timeout... shutting down")
		s.Shutdown()
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.serverPort))
	if err != nil {
		log.Println(err)
		return
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
		s.client = conn
		conn.Write([]byte("Welcome to this "))
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
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

func (s *Server) HandleFile(w http.ResponseWriter, r *http.Request) {
	// upgrade connectin to websocket
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn1, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn1.Close()

	s.logger.WriteLog("sending file")

	/*
		Get client connection and put it here
	*/

	// Create a channel to forward messages from conn1 to conn2
	forwardChan := make(chan []byte)

	// Start a goroutine to read messages from conn1 and forward them to chan
	go func() {
		for {
			messageType, message, err := conn1.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}
			if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
				forwardChan <- message
			}
		}
	}()

	// Start a goroutine to read messages from forwardChan and write them to conn2
	go func() {
		for message := range forwardChan {
			i, err := s.client.Write(message)
			if err != nil {
				s.logger.WriteLog(fmt.Sprint("Error writing message:", err))
				break
			}
			s.logger.WriteLog(fmt.Sprintf("Written %d bytes", i))
		}
	}()

	// Wait for either goroutine to exit
	select {}
}

func (s *Server) Send(filePath string) error {
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
		err = binary.Write(s.client, binary.LittleEndian, int64(len(data)))
		if err != nil {
			s.logger.WriteLog(err)
			return
		}

		// send data
		i, err := io.CopyN(s.client, bytes.NewReader(data), int64(len(data)))
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
