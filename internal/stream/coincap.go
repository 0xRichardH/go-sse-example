package stream

import (
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"richard/go-sse-demo/internal/sse"
)

type CoinCap struct {
	url    url.URL
	broker *sse.Broker
}

func NewCoinCapClient(broker *sse.Broker) (c *CoinCap) {
	// u := url.URL{Scheme: "wss", Host: "ws.coincap.io", Path: "prices", RawQuery: "assets=bitcoin,ethereum,monero,litecoin"}
	u := url.URL{Scheme: "wss", Host: "ws.coincap.io", Path: "trades/binance"}
	log.Printf("connecting to %s", u.String())

	c = &CoinCap{
		url:    u,
		broker: broker,
	}

	c.dial()

	return
}

func (c *CoinCap) dial() {
	client, _, err := websocket.DefaultDialer.Dial(c.url.String(), nil)
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
			// log.Printf("recv: %s", message)
			c.broker.Notifier <- message
		}
	}()

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

			<-done
			return
		}
	}
}
