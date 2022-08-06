package stream

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"time"
)

type Streamer struct {
}

func NewStreamer() (streamer *Streamer) {
	streamer = &Streamer{}
	return
}

func (s *Streamer) ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("upgrade:", err)
	}
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		err = writeMessage(conn, string(message))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func writeMessage(conn *websocket.Conn, message string) (err error) {
	var writeMessage string
	switch message {
	case "ping":
		writeMessage = "pong"
	case "time":
		writeMessage = time.Now().String()
	default:
		writeMessage = "Do I know you?"
	}
	for {
		time.Sleep(time.Second * 2)
		err = conn.WriteMessage(websocket.TextMessage, []byte(writeMessage))
		if err != nil {
			break
		}
	}
	return
}