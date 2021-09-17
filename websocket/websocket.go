// Package websocket contains interfaces to represent websocket structures
package websocket

import (
	"io"
	"net/http"
)

// Upgrader is used to upgrade an existing connection to a websocket connection
type Upgrader interface {
	// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
	Upgrade(http.ResponseWriter, *http.Request, http.Header) (WebsocketConnection, error)
}

// WebsocketConnection represents a single websocket connection
type WebsocketConnection interface {
	// Close closes the websocket connection
	Close() error

	// NextWriter returns a writer for the next message to send
	NextWriter(messageType int) (io.WriteCloser, error)
}
