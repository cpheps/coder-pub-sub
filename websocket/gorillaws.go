package websocket

import (
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

	return &GorillaConn{conn}, nil
}

var _ (WebsocketConnection) = (*GorillaConn)(nil)

// GorillaConn is a wrapper around the gorilla/websocket Conn to satisfy the WebsocketConnection interface
type GorillaConn struct {
	*gwebsocket.Conn
}
