package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/cpheps/coder-pub-sub/server"
)

func main() {
	// Setup signal context
	signalCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	pubsubServer, err := server.New(":8080", 10)
	if err != nil {
		log.Fatalln("Failed to init server", err)
	}

	// Spin server off in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- pubsubServer.ListenAndServe()
	}()

	// Wait for a single to close or the server to exit
	select {
	case <-signalCtx.Done():
		if err := pubsubServer.Close(); err != nil {
			log.Fatalln("Error while closing server", err)
		}
	case err := <-errChan:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("Server closed with error", err)
		}
	}
}
