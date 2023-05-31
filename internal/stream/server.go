package stream

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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
			log.Fatal("read:", err)
			return
		}
		log.Printf("recv: %s", message)

		for {
			time.Sleep(time.Second * 2)
			err = writeMessage(conn, string(message))
			if err != nil {
				log.Fatal("write:", err)
				return
			}
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
	err = conn.WriteMessage(websocket.TextMessage, []byte(writeMessage))
	return
}
