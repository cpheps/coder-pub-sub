// Package websocket contains interfaces to represent websocket structures
package websocket

import (
	"io"
	"net/http"
)

// MessageType is defined as in RFC 6455 https://datatracker.ietf.org/doc/html/rfc6455#section-11.8
type MessageType int

const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage MessageType = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage MessageType = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage MessageType = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage MessageType = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage MessageType = 10
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
	NextWriter(messageType MessageType) (io.WriteCloser, error)
}
