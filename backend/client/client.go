package client

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"swift2/global"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn           net.Conn
	hostname       string
	AvailableHosts []Host
	logger         *global.Logger
}

func NewClient() *Client {
	return &Client{
		hostname: fmt.Sprint("", rand.Intn(20)),
	}
}

type Host struct {
	hostname string
	ipPort   string
}

func (c *Client) Listen() {
	// Resolve the broadcast address and port
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", global.BroadcastPort))
	if err != nil {
		return
	}

	// Create a UDP socket to listen on
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return
	}
	defer conn.Close()

	// Set a timeout for the socket
	conn.SetReadDeadline(time.Now().Add(time.Second * 30))

	println("now listening on port ", global.BroadcastPort)

	// Wait for a message
	stopTime := time.Now().Add(time.Second * 7)
	for time.Now().Before(stopTime) {
		buffer := new(bytes.Buffer)
		_, _, err := conn.ReadFromUDP(buffer.Bytes())
		// fmt.Println("broadcast received: ", n, remoteAddr, buffer.String())
		// c.AvailableHosts = append(c.AvailableHosts, Host{hostname: "", ipPort: })
		if err != nil {
			return
		}
	}

}

func (c *Client) Connect(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}

	c.conn = conn

	// send hostname
	go c.readLoop(conn)

}

func (c *Client) readLoop(conn net.Conn) {
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

func (c *Client) HandleFile(w http.ResponseWriter, r *http.Request) {
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
		err = binary.Write(c.conn, binary.LittleEndian, int64(len(jsonMsg)))
		if err != nil {
			log.Println(err)
			return
		}
		io.CopyN(c.conn, bytes.NewReader(jsonMsg), int64(len(jsonMsg)))
		c.logger.WriteLog(fmt.Sprint("transfer complete: ", len(jsonMsg), "bytes sent"))
	}
}
func (c *Client) Send(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return err
	}

	// prefix data with filename
	filename := []byte(fmt.Sprintf("%v$$$$", path.Base(filePath)))
	data = append(filename, data...)

	// send file size
	err = binary.Write(c.conn, binary.LittleEndian, int64(len(data)))
	if err != nil {
		log.Println(err)
		return err
	}

	// send data
	i, err := io.CopyN(c.conn, bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("File sent successfully: %d / %d bytes written", i, len(data))

	return nil
}

func (c *Client) Receive() {

	// receive file size
	var dataSize int64
	err := binary.Read(c.conn, binary.LittleEndian, &dataSize)
	if err != nil {
		log.Fatal(err)
	}

	// recieve data prefixed with filename
	data := new(bytes.Buffer)
	i, err := io.CopyN(data, c.conn, dataSize)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Received %d bytes from connection", i)

	// seperate data from filename
	filename, fileContent, ok := bytes.Cut(data.Bytes(), []byte("$$$$"))
	if !ok {
		log.Println("Unable to parse file... ")
	}

	err = os.WriteFile(fmt.Sprintf("./%s%s", c.hostname, filename), fileContent, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) Disconnect() error {
	err := c.conn.Close()
	if err != nil {
		log.Println(err)
	}
	return err
}
