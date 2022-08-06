package stream

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type CoinCap struct {
	Url         url.URL
	TradeStream chan []byte
}

func NewCoinCapClient() (c *CoinCap) {
	u := url.URL{Scheme: "wss", Host: "ws.coincap.io", Path: "prices", RawQuery: "assets=bitcoin,ethereum,monero,litecoin"}
	// u := url.URL{Scheme: "wss", Host: "ws.coincap.io", Path: "trades/binance"}
	log.Printf("connecting to %s", u.String())

	c = &CoinCap{
		Url:         u,
		TradeStream: make(chan []byte, 1),
	}

	c.dial()

	return
}

func (c *CoinCap) dial() {
	client, response, err := websocket.DefaultDialer.Dial(c.Url.String(), nil)
	if err != nil {
		if response != nil {
			log.Printf("handshake failed with status %d", response.StatusCode)
		} else {
			log.Printf("handshake failed without response")
		}
		log.Fatal("dial:", err)
	}

	defer client.Close()

	go func() {
		for {
			_, message, err := client.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			c.TradeStream <- message
		}
	}()
}
