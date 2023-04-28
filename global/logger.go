package global

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Logger struct {
	conn websocket.Conn
}

func NewLogger(socket websocket.Conn) *Logger {
	return &Logger{
		conn: socket,
	}
}

func (l *Logger) Write(logType string, msg any) {
	l.conn.WriteJSON(fmt.Sprintf(`{"%s": "%s"}`, logType, fmt.Sprint(msg)))
	log.Println(msg)
}
