package stream

import (
	"context"
	"log"
	"net/url"

	"richard/go-sse-demo/internal/sse"

	"github.com/gorilla/websocket"
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

	return
}

func (s *Stream) Dial(ctx context.Context) {
	client, _, err := websocket.DefaultDialer.Dial(s.url.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer client.Close()

	ctx, cancelCtx := context.WithCancel(ctx)
	go func() {
		defer cancelCtx()
		for {
			_, message, err := client.ReadMessage()
			if err != nil {
				log.Fatal("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			s.broker.Notifier <- message
		}
	}()

	requestTimeStream(client)

	<-ctx.Done()
	if err := ctx.Err(); err != nil {
		log.Fatal("ctx:", err)
	}
	err = client.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Fatal("write close:", err)
	}
}

func requestTimeStream(c *websocket.Conn) {
	err := c.WriteMessage(websocket.TextMessage, []byte("time"))
	if err != nil {
		log.Fatal("write:", err)
	}
}
