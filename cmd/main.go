package main

import (
	"log"
	"net/http"

	"richard/go-sse-demo/internal/sse"
	"richard/go-sse-demo/internal/stream"
)

func main() {
	broker := sse.NewServer()
	// coincap := stream.NewCoinCapClient()
	streamClient := stream.NewStreamClient()

	// go func() {
	// 	for {
	// 		time.Sleep(2 * time.Second)
	// 		eventString := fmt.Sprintf("the time is %v", time.Now())
	// 		log.Println("Receiving event")
	// 		broker.Notifier <- []byte(eventString)
	// 	}
	// }()

	// go func() {
	// 	message := <-coincap.TradeStream
	// 	log.Println("Receiving trade stream.")
	// 	broker.Notifier <- message
	// }()

	go func() {
		message := <-streamClient.MessageChan
		broker.Notifier <- message
	}()

	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:3000", broker))

}
