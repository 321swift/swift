package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

type UiServer struct {
	socket *websocket.Conn
	port   int
	role   string
}

func main() {
	// start server
	server := NewUiServer()
	server.Start()
}
func NewUiServer() *UiServer {
	return &UiServer{
		port: 3000,
	}
}

func OpenPage(url string) {
	var err error

	switch runtime.GOOS {
	case "darwin":
		err = exec.Command("open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start()
	}

	if err != nil {
		fmt.Println("Failed to open URL:", err)
	}

}

func (s *UiServer) Start() {
	fmt.Println("starting uiserver")
	for {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
		if err == nil {
			r := chi.NewRouter()
			r.Use(middleware.Logger)

			r.Handle("/", http.FileServer(http.Dir("./ui")))
			r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static"))))
			r.HandleFunc("/ws", s.HandleWS)
			defer listener.Close()

			fmt.Printf("Server started on port %d\n", s.port)
			OpenPage(fmt.Sprintf("http://localhost:%d", s.port))
			http.Serve(listener, r)

			return
		}
		// Port is already in use, so try the next one
		fmt.Printf("Port %d already in use, trying next port\n", s.port)
		s.port++
	}
}

func (s *UiServer) HandleWS(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	s.socket = conn
	conn.WriteJSON("status: ok")

	go s.ReadLoop(w)
}

func (s *UiServer) ReadLoop(w http.ResponseWriter) {
	defer s.socket.Close()
	for {
		_, content, err := s.socket.ReadMessage()
		if err != nil {
			log.Println(err)
			s.socket.WriteJSON("server err")
		}
		log.Println(string(content))

		var roleStruct = &struct{ Role string }{}
		err = json.Unmarshal(content, roleStruct)
		if err != nil {
			log.Println("unable to parse json:", err)
			return
		}
		switch roleStruct.Role {
		case "server":
			s.socket.WriteJSON("assuming server role")

		case "client":
			s.socket.WriteJSON("assuming client role")
		}
	}
}

func assumeServer() {

}
func assumeClient() {

}
