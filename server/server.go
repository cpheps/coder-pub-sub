// Package server contains the implementation of the PubSub server
package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/cpheps/coder-pub-sub/websocket"
	"github.com/gorilla/mux"
	gwebsocket "github.com/gorilla/websocket"
)

// PubSubServer a server the implements pub/sub via websockets
type PubSubServer struct {
	doneChan    chan struct{}
	srv         *http.Server
	upgrader    websocket.Upgrader
	broadcaster websocket.Broadcaster
}

// New creates a new instance of the PubSub Server that listens on the supplied addr.
// Uses gorilla websocket and mux
func New(addr string, broadcastConcurrency int) (*PubSubServer, error) {
	broadcaster, err := websocket.NewCacheBroadcaster(broadcastConcurrency)
	if err != nil {
		return nil, err
	}

	// Create the server before the router so we can register it's handlers on the router
	pubSubServer := &PubSubServer{
		doneChan: make(chan struct{}),
		srv: &http.Server{
			Addr: addr,
		},
		upgrader:    websocket.NewGorillaUpgrader(&gwebsocket.Upgrader{}),
		broadcaster: broadcaster,
	}

	r := mux.NewRouter()

	// Register GET only for subscribe
	r.HandleFunc("/subscribe", pubSubServer.RegisterSubscriber).Methods(http.MethodGet)

	// Register Post only for publish
	r.HandleFunc("/publish", pubSubServer.Publish).Methods(http.MethodPost)

	// Set mux on the server
	pubSubServer.srv.Handler = r
	return pubSubServer, nil
}

// ListenAndServe starts the server and blocks until the server returns.
// The server can be closed via the Close call
func (s *PubSubServer) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

// Close causes a graceful shutdown of the server
func (s *PubSubServer) Close() error {
	// Close the done channel to stop all blocking handlers
	close(s.doneChan)
	s.broadcaster.CloseConnections()
	return s.srv.Close()
}

// RegisterSubscriber registers a subscriber with the server and opens up a websocket
func (s *PubSubServer) RegisterSubscriber(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrading connection to websocket", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Register the connection
	s.broadcaster.RegisterConnection(conn)

	// Block until server closes as we don't want the websocket to prematurely die
	<-s.doneChan
}

// Publish publishes a messsage to all subscribers
func (s *PubSubServer) Publish(w http.ResponseWriter, r *http.Request) {
	// Parse the message body
	defer r.Body.Close()
	msg, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error while reading message body", err)
		s.writeResponse(w, http.StatusBadRequest, &errorResponse{
			Message: "failed to read message",
		})
		return
	}

	// Broadcast message
	// Hard coding to messageType of TextMessag but ideally could parse the ContentType header and dynamically change
	if err := s.broadcaster.Broadcast(r.Context(), websocket.TextMessage, msg); err != nil {
		log.Println("Broadcase failure", err)
		s.writeResponse(w, http.StatusInternalServerError, &errorResponse{
			Message: "Internal Error",
		})
	}

	// Success no content
	w.WriteHeader(http.StatusNoContent)
}

func (s *PubSubServer) writeResponse(w http.ResponseWriter, code int, v interface{}) {
	w.WriteHeader(code)

	payload, err := json.Marshal(v)
	if err != nil {
		log.Println("failed to marshal payload:", err)
	}

	_, err = w.Write(payload)
	if err != nil {
		log.Println("failed to write payload:", err)
	}
}
