package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"path"
	"time"

	"github.com/321swift/swift/global"
)

type client struct {
	conn           net.Conn
	hostname       string
	AvailableHosts []Host
}

func NewClient() *client {
	return &client{
		hostname: fmt.Sprint("", rand.Intn(20)),
	}
}

type Host struct {
	hostname string
	ipPort   string
}

func (c *client) Listen() {
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

func (c *client) Connect(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	c.conn = conn

	// send hostname
	go c.readLoop(conn)

}

func (c *client) readLoop(conn net.Conn) {
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
		log.Printf("Received %d bytes from %v", i, conn.RemoteAddr())

		// seperate data from filename
		filename, fileContent, ok := bytes.Cut(data.Bytes(), []byte("$$$$"))
		if !ok {
			log.Println("Unable to parse file... ")
		}

		// err = os.WriteFile(fmt.Sprintf("./1%s", filename), fileContent, os.ModePerm)
		err = os.WriteFile(fmt.Sprintf("./%s%s", c.hostname, filename), fileContent, os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}
	}

}

func (c *client) Send(filePath string) error {
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

func (c *client) Receive() {

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

func (c *client) Disconnect() error {
	err := c.conn.Close()
	if err != nil {
		log.Println(err)
	}
	return err
}
