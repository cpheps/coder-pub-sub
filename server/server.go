// Package server contains the implementation of the PubSub server
package server

import (
	"net/http"

	"github.com/cpheps/coder-pub-sub/websocket"
	"github.com/gorilla/mux"
	gwebsocket "github.com/gorilla/websocket"
)

// PubSubServer a server the implements pub/sub via websockets
type PubSubServer struct {
	srv      *http.Server
	upgrader websocket.Upgrader
}

// New creates a new instance of the PubSub Server that listens on the supplied addr.
// Uses gorilla websocket and mux
func New(addr string) *PubSubServer {
	// Create the server before the router so we can register it's handlers on the router
	pubSubServer := &PubSubServer{
		srv: &http.Server{
			Addr: addr,
		},
		upgrader: websocket.NewGorillaUpgrader(&gwebsocket.Upgrader{}),
	}

	r := mux.NewRouter()

	// Register GET only for subscribe
	r.HandleFunc("/subscribe", pubSubServer.RegisterSubscriber).Methods(http.MethodGet)

	// Register Post only for publish
	r.HandleFunc("/publish", pubSubServer.Publish).Methods(http.MethodPost)

	// Set mux on the server
	pubSubServer.srv.Handler = r
	return pubSubServer
}

// ListenAndServe starts the server and blocks until the server returns.
// The server can be closed via the Close call
func (s *PubSubServer) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

// Close causes a graceful shutdown of the server
func (s *PubSubServer) Close() error {
	return s.srv.Close()
}

// RegisterSubscriber registers a subscriber with the server and opens up a websocket
func (s *PubSubServer) RegisterSubscriber(w http.ResponseWriter, r *http.Request) {

}

// Publish publishes a messsage to all subscribers
func (s *PubSubServer) Publish(w http.ResponseWriter, r *http.Request) {

}
