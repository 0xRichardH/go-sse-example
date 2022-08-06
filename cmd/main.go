package main

import (
	"log"
	"net/http"

	"richard/go-sse-demo/internal/sse"
	"richard/go-sse-demo/internal/stream"
)

func main() {
	broker := sse.NewServer()

	// go stream.NewCoinCapClient(broker)
	go stream.NewStreamClient(broker)

	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:8081", broker))
}
