package node

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"swift/global"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
)

type Node struct {
	basePort          int
	infoLog           *log.Logger
	errLog            *log.Logger
	connectionPool    map[int]net.Listener
	uiSocket          *websocket.Conn
	senderTimeout     time.Duration
	serverPort        int
	hostname          string
	backendConnection net.Conn
}
type Intro struct {
	hostname string
	conns    []int
}

func NewNode(infoLog *log.Logger, errLog *log.Logger) *Node {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = fmt.Sprintf("swift%d", rand.Intn(500))
	}

	return &Node{
		basePort:       4009,
		infoLog:        infoLog,
		errLog:         errLog,
		connectionPool: make(map[int]net.Listener),
		senderTimeout:  time.Second * 20,
		hostname:       hostname,
	}
}

func (n *Node) Start() {
	// setup connection pool
	for i := 0; i < 5; i++ {
		port := n.getAvailablePort()
		n.connectionPool[port] = nil
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			n.errLog.Println(err)
			return
		}
		n.connectionPool[port] = listener
		go func() {
			fileRouter := chi.NewRouter()
			fileRouter.Use(middleware.Logger)
			fileRouter.HandleFunc("/file", n.handleFileReception)

			http.Serve(listener, fileRouter)
		}()
	}
	fmt.Printf("%+v\n", n.connectionPool)

	// begin UI server
	func() {
		n.infoLog.Println("Starting UI server")
		port := n.getAvailablePort()
		r := chi.NewRouter()
		r.Use(middleware.Logger)

		r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static"))))
		r.Handle("/*", http.FileServer(http.Dir("ui/")))
		r.HandleFunc("/sender", n.handleSenderRole)
		r.HandleFunc("/receiver", n.handleReceiverRole)
		n.infoLog.Printf("Server started on port %d\n", port)
		OpenPage(fmt.Sprintf("http://localhost:%d", port))
		http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	}()
}

func (n *Node) handleSenderRole(w http.ResponseWriter, r *http.Request) {
	// server role assumed
	// upgrade to websocket
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	n.uiSocket = conn
	conn.WriteJSON(`{"status": "sender"}`)
	n.infoLog.Println(`{"status": "sender"}`)

	go n.ReadLoop(n.uiSocket)()

	// broadcast
	n.serverPort = n.getAvailablePort()
	go func() {
		n.broadcast()
	}()

	n.infoLog.Println("backend server started on port ", n.serverPort)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", n.serverPort))
	if err != nil {
		n.errLog.Println(listener)
		return
	}
	timer := time.AfterFunc(n.senderTimeout, func() {
		n.infoLog.Println("status", "Server timeout... shutting down")
		// close backend listener
		listener.Close()
	})

	n.backendConnection, err = listener.Accept()
	if err != nil {
		n.errLog.Println("accepting connection err: ", err)
	}
	timer.Stop()

	if n.backendConnection != nil {
		n.senderIntroduce()
		n.ReadLoop(n.backendConnection)()
	} else {
		n.uiSocket.WriteJSON(`{"status":"server timed out; no connections made"}`)
	}
}

func (n *Node) senderIntroduce() {
	// extract ports of connection pool
	keys := make([]int, 0)
	for i := range n.connectionPool {
		keys = append(keys, i)
	}

	// send intro
	intro := Intro{
		hostname: n.hostname,
		conns:    keys,
	}
	bin_buf := new(bytes.Buffer)
	gobObj := gob.NewEncoder(bin_buf)
	gobObj.Encode(intro)
	_, err := n.backendConnection.Write(bin_buf.Bytes())
	if err != nil {
		n.errLog.Println(err)
		return
	}

	// read intro
	tmp := make([]byte, 500)
	_, err = n.backendConnection.Read(tmp)
	if err != nil {
		n.errLog.Println(err)
	}

	tmpbuff := bytes.NewBuffer(tmp)
	receivedIntro := new(Intro)

	gobDec := gob.NewDecoder(tmpbuff)

	gobDec.Decode(receivedIntro)
	n.infoLog.Println(receivedIntro)
}

func (n *Node) receiverIntroduce() {

	keys := make([]int, 0)
	for i := range n.connectionPool {
		keys = append(keys, i)
	}

	// read intro
	tmp := make([]byte, 500)
	_, err := n.backendConnection.Read(tmp)
	if err != nil {
		n.errLog.Println(err)
	}

	tmpbuff := bytes.NewBuffer(tmp)
	receivedIntro := new(Intro)

	gobDec := gob.NewDecoder(tmpbuff)
	gobDec.Decode(receivedIntro)

	// send intro
	intro := Intro{
		hostname: n.hostname,
		conns:    keys,
	}
	bin_buf := new(bytes.Buffer)
	gobObj := gob.NewEncoder(bin_buf)
	gobObj.Encode(intro)
	n.backendConnection.Write(bin_buf.Bytes())

	n.infoLog.Println(receivedIntro)
}

func (n *Node) handleReceiverRole(w http.ResponseWriter, r *http.Request) {
	// receiver role assumed
	// upgrade to websocket
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		n.errLog.Println(err)
		return
	}
	n.uiSocket = conn
	conn.WriteJSON(`{"status": "receiver"}`)
	n.infoLog.Println(`{"status": "receiver"}`)

	go n.ReadLoop(n.uiSocket)()

	// start listening
	availableHost, err := n.Listen()
	if err != nil {
		n.errLog.Println(err)
		return
	}
	err = n.Connect(availableHost)
	if err != nil {
		n.errLog.Println(err)
		return
	}
	n.receiverIntroduce()

}

func (n *Node) Listen() (string, error) {
	// Resolve the broadcast address and port
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", BroadcastPort))
	if err != nil {
		return "", err
	}

	// Create a UDP socket to listen on
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// Set a timeout for the socket
	conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	n.infoLog.Println("now listening on port ", global.BroadcastPort)

	// Wait for a message
	stopTime := time.Now().Add(n.senderTimeout)
	availableHost := ""
	for time.Now().Before(stopTime) {
		buffer := make([]byte, 40)
		_, remoteAddr, err := conn.ReadFromUDP(buffer)
		// c.AvailableHosts = append(c.AvailableHosts, Host{hostname: "", ipPort: })
		if err != nil {
			return "", err
		}

		n.infoLog.Println("broadcast received: ", remoteAddr.String(), string(buffer))
		availableHost = fmt.Sprintf("%s:%s",
			strings.Split(remoteAddr.String(), ":")[0],
			strings.Split(string(buffer), ":")[1],
		)
		if availableHost != "" {
			break
		}
	}
	return availableHost, nil
}

func (n *Node) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	n.backendConnection = conn

	// send hostname
	return nil

}

func (n *Node) ReadLoop(conn interface{}) func() {
	switch v := conn.(type) {
	case net.Conn:
		return func() {
			for {
				readBytes := new(bytes.Buffer)
				_, err := v.Read(readBytes.Bytes())
				if err != nil {
					n.errLog.Println(err)
					break
				}

				n.infoLog.Println(readBytes.String())
			}

		}
	case *websocket.Conn:
		return func() {
			for {
				_, content, err := v.ReadMessage()
				if err != nil {
					n.errLog.Println(err)
					break
				}

				n.infoLog.Println(content)
			}

		}
	default:
		return nil
	}
}
