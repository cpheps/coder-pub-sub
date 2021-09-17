package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/cpheps/coder-pub-sub/server"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  0,
	WriteBufferSize: 1024,
}

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
