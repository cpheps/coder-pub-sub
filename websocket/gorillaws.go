package websocket

import (
	"io"
	"net/http"

	gwebsocket "github.com/gorilla/websocket"
)

var _ (Upgrader) = (*GorillaUpgrader)(nil)

// GorillaUpgrader is a wrapper around the gorilla/websocket Upgrader to satisfy the Upgrader interface
type GorillaUpgrader struct {
	upgrader *gwebsocket.Upgrader
}

// NewGorillaUpgrader creates a new GorillaUpgrader that wraps the passed in upgrader
func NewGorillaUpgrader(upgrader *gwebsocket.Upgrader) *GorillaUpgrader {
	return &GorillaUpgrader{
		upgrader: upgrader,
	}
}

// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
func (gu *GorillaUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (WebsocketConnection, error) {
	conn, err := gu.upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	return &GorillaConn{
		conn: conn,
	}, nil
}

var _ (WebsocketConnection) = (*GorillaConn)(nil)

// GorillaConn is a wrapper around the gorilla/websocket Conn to satisfy the WebsocketConnection interface
type GorillaConn struct {
	conn *gwebsocket.Conn
}

// NextWriter returns a writer for the next message to send
func (gc *GorillaConn) NextWriter(messageType MessageType) (io.WriteCloser, error) {
	return gc.conn.NextWriter(int(messageType))
}

// Close closes the websocket connection
func (gc *GorillaConn) Close() error {
	return gc.conn.Close()
}
