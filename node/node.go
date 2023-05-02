package node

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
)

type Node struct {
	basePort       int
	infoLog        *log.Logger
	errLog         *log.Logger
	connectionPool map[int]net.Listener
	uiSocket       *websocket.Conn
}

func NewNode(infoLog *log.Logger, errLog *log.Logger) *Node {
	return &Node{
		basePort:       5050,
		infoLog:        infoLog,
		errLog:         errLog,
		connectionPool: make(map[int]net.Listener),
	}
}

func (n *Node) Start() {
	// setup connection pool
	for i := 0; i < 5; i++ {
		port := n.GetAvailablePort()
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
		port := n.GetAvailablePort()
		r := chi.NewRouter()
		r.Use(middleware.Logger)

		r.Handle("/*", http.FileServer(http.Dir("ui/")))
		r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static"))))
		r.HandleFunc("/sender", n.handleServerRole)
		r.HandleFunc("/receiver", n.handleReceiverRole)
		n.infoLog.Printf("Server started on port %d\n", port)
		OpenPage(fmt.Sprintf("http://localhost:%d", port))
		http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	}()
}

func (n *Node) handleServerRole(w http.ResponseWriter, r *http.Request) {
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
	defer conn.Close()
	n.uiSocket = conn
	conn.WriteJSON(`{"status": "sender"}`)

	n.ReadLoop()
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
		log.Println(err)
		return
	}
	defer conn.Close()
	n.uiSocket = conn
	conn.WriteJSON(`{"status": "receiver"}`)

	n.ReadLoop()
}

func (n *Node) ReadLoop() {
	defer n.uiSocket.Close()
	for {
		_, content, err := n.uiSocket.ReadMessage()
		if err != nil {
			n.errLog.Println(err)
		}

		n.infoLog.Println(content)
	}
}
