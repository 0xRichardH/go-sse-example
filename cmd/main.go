package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"richard/go-sse-demo/internal/sse"
	"richard/go-sse-demo/internal/stream"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	broker := sse.NewServer(ctx)
	server := http.Server{
		Addr:    ":3000",
		Handler: broker,
	}

	// go stream.NewCoinCapClient(broker)
	streamClient := stream.NewStreamClient(broker)
	go streamClient.Dial(ctx)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	// Perform application shutdown with a maximum timeout of 5 seconds.
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}
