package stream

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Stream struct {
	url         url.URL
	MessageChan chan []byte
}

func NewStreamClient() (s *Stream) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "bot"}
	log.Printf("connecting to %s", u.String())

	s = &Stream{
		url:         u,
		MessageChan: make(chan []byte, 1),
	}

	s.dial()

	return
}

func (s *Stream) dial() {
	client, _, err := websocket.DefaultDialer.Dial(s.url.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer client.Close()

	go func() {
		for {
			_, message, err := client.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)
			s.MessageChan <- message
		}
	}()

	requestTimeStream(client)
}

func requestTimeStream(c *websocket.Conn) {
	err := c.WriteMessage(websocket.TextMessage, []byte("time"))
	if err != nil {
		log.Println("write:", err)
	}
}
