package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"richard/go-sse-demo/internal/sse"
)

func main() {
	broker := sse.NewServer()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			eventString := fmt.Sprintf("the time is %v", time.Now())
			log.Println("Receiving event")
			broker.Notifier <- []byte(eventString)
		}
	}()

	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:3000", broker))
}
