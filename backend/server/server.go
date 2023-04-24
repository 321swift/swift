package server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	log.Println("Receiving file from connection")

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
		log.Printf("Received %d bytes from %v", i, conn.RemoteAddr())

		// Decode the JSON-encoded message and extract the filename and data
		var msg global.Message
		err = json.Unmarshal(data.Bytes(), &msg)
		if err != nil {
			log.Println("Error decoding message:", err)
			break
		}

		// Write the data to a file with the given filename
		err = ioutil.WriteFile(msg.Filename, msg.Data, 0644)
		if err != nil {
			log.Println("Error writing file:", err)
			break
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
	messageType, message, err := conn1.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
	}
	if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
		// Create a Message struct and encode it as JSON
		msg := global.Message{}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("Error encoding message:", err)
		}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Println("Error encoding message:", err)
		}
		err = binary.Write(s.client, binary.LittleEndian, int64(len(jsonMsg)))
		if err != nil {
			log.Println(err)
			return
		}
		io.CopyN(s.client, bytes.NewReader(jsonMsg), int64(len(jsonMsg)))
		s.logger.WriteLog(fmt.Sprint("transfer complete: ", len(jsonMsg), "bytes sent"))
	}
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
