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
	streamClient := stream.NewStreamClient(broker)
	go streamClient.Dial()

	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:8088", broker))
}
