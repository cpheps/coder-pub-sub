package websocket

import (
	"context"
	"errors"
	"fmt"
	"log"

	"golang.org/x/sync/errgroup"
)

// Broadcaster broadcasts a message to several websockets
type Broadcaster interface {
	// RegisterConnection registers a connection with the Broadcaster
	RegisterConnection(WebsocketConnection)

	// Broadcast sends the bytes of messageType to all websockets.
	// Returns and error if a single send fails
	Broadcast(ctx context.Context, messageType MessageType, msg []byte) error

	// CloseConnections closes all registered connections
	CloseConnections()
}

var _ (Broadcaster) = (*CacheBroadcaster)(nil)

// CacheBroadcaster implements the Broadcaster interface as well as locally caches websocket connections
type CacheBroadcaster struct {
	conns []WebsocketConnection

	// concurrency is the number of goroutines to have active at a time while sending
	concurrency int
}

// NewCacheBroadcaster creates a new CacheBroadcaster with the passed in concurrency
func NewCacheBroadcaster(concurrency int) (*CacheBroadcaster, error) {
	if concurrency <= 0 {
		return nil, errors.New("concurrency must be greater than 0")
	}

	return &CacheBroadcaster{
		conns:       make([]WebsocketConnection, 0),
		concurrency: concurrency,
	}, nil
}

// RegisterConnection registers a connection with the Broadcaster
func (cb *CacheBroadcaster) RegisterConnection(conn WebsocketConnection) {
	cb.conns = append(cb.conns, conn)
}

// CloseConnections closes all registered connections
// Will log any errors
func (cb *CacheBroadcaster) CloseConnections() {
	for _, conn := range cb.conns {
		if err := conn.Close(); err != nil {
			log.Println("Error while closing websocket", err)
		}
	}

	// After all connections are closed clean our connection tracking
	cb.conns = make([]WebsocketConnection, 0)
}

// Broadcast sends the bytes of messageType to all websockets.
// Returns and error if a single send fails
func (cb *CacheBroadcaster) Broadcast(ctx context.Context, messageType MessageType, msg []byte) error {
	group, errCtx := errgroup.WithContext(ctx)

	// Create a buffered channel large enough so each worker is busy
	socketChan := make(chan WebsocketConnection, cb.concurrency)

	// Spin up workers to handle broadcasting
	for i := 0; i < cb.concurrency; i++ {
		group.Go(func() error {
			return broadcastWorker(errCtx, socketChan, messageType, msg)
		})
	}

	// Feed connections to workers
	for _, conn := range cb.conns {
		socketChan <- conn
	}

	close(socketChan)

	return group.Wait()
}

// broadcastWorker sends the message to each WebsocketConnection supplied to it
func broadcastWorker(ctx context.Context, socketChan <-chan WebsocketConnection, messageType MessageType, msg []byte) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case conn, ok := <-socketChan:
			if !ok {
				return nil
			}

			// Create a new writer for the websocket
			writer, err := conn.NextWriter(messageType)
			if err != nil {
				return fmt.Errorf("failed to create writer for websocket: %w", err)
			}

			// Write all data to the writer created by the connection
			for written := 0; written < len(msg); {
				numBytes, err := writer.Write(msg[written:])
				if err != nil {
					return fmt.Errorf("failed while writing to websocket: %w", err)
				}

				written += numBytes
			}

			// Close the writer
			if err := writer.Close(); err != nil {
				return fmt.Errorf("failed to close writer for websocket: %w", err)
			}
		}
	}
}
