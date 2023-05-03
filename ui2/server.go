package ui2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"swift/backend/server"
	"swift/global"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
)

type UiServer struct {
	socket        *websocket.Conn
	port          int
	filePort      int
	fskt          bool
	backendServer *server.Server
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

func (s *UiServer) StartFileSocket() {
	fmt.Println("starting uiserver")
	s.filePort = global.GetAvailablePort(s.port)
	r := chi.NewRouter()

	r.HandleFunc("/file", s.backendServer.HandleFile)

	fmt.Printf("file Server started on port %d\n", s.filePort)
	s.socket.WriteJSON(fmt.Sprintf(`{"fileSocket": %d}`, s.filePort))
	http.ListenAndServe(fmt.Sprintf(":%d", s.filePort), r)
}

func (s *UiServer) Start() {
	fmt.Println("starting uiserver")
	s.port = global.GetAvailablePort(3000)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Handle("/*", http.FileServer(http.Dir("ui2/")))
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("ui2/static"))))
	r.HandleFunc("/ws", s.HandleWS)
	fmt.Printf("Server started on port %d\n", s.port)
	OpenPage(fmt.Sprintf("http://localhost:%d", s.port))
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), r)
}

func (s *UiServer) HandleWS(w http.ResponseWriter, r *http.Request) {
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
	s.socket = conn
	conn.WriteJSON(`{"status": "ok"}`)
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
		// log.Println(string(content))

		var roleStruct = &struct{ Role string }{}
		err = json.Unmarshal(content, roleStruct)
		if err != nil {
			log.Println("unable to parse json:", err)
			return
		}
		switch roleStruct.Role {
		case "server":
			s.socket.WriteJSON(`{"status": "assuming server role"}`)
			s.assumeServer(*s.socket)
		case "client":
			s.socket.WriteJSON("assuming client role")
		}
	}
}
func (s *UiServer) assumeServer(conn websocket.Conn) {
	defer conn.Close()
	logger := global.NewLogger(conn)

	s.backendServer = server.NewServer(logger)

	if !s.fskt {
		s.fskt = true
		go s.StartFileSocket()
	} else {
		s.socket.WriteJSON(fmt.Sprintf(`{"fileSocket": %d}`, s.filePort))
	}

	s.backendServer.Start()
}
