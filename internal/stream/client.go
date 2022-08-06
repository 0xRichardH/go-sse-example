package stream

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"os"
	"os/signal"
	"richard/go-sse-demo/internal/sse"
)

type Stream struct {
	url    url.URL
	broker *sse.Broker
}

func NewStreamClient(broker *sse.Broker) (s *Stream) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "bot"}
	log.Printf("connecting to %s", u.String())

	s = &Stream{
		url:    u,
		broker: broker,
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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := client.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			s.broker.Notifier <- message
		}
	}()

	requestTimeStream(client)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")
			err := client.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}

			select {
			case <-done:
			}

			return
		}
	}
}

func requestTimeStream(c *websocket.Conn) {
	err := c.WriteMessage(websocket.TextMessage, []byte("time"))
	if err != nil {
		log.Println("write:", err)
	}
}
