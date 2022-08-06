package main

import (
	"log"
	"net/http"
	"richard/go-sse-demo/internal/stream"
)

func main() {
	streamer := stream.NewStreamer()
	http.HandleFunc("/bot", streamer.ServeWebSocket)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
