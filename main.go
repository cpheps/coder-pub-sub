package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  0,
	WriteBufferSize: 1024,
}

func main() {

	http.HandleFunc("/", WebSocket)
	if err := http.ListenAndServe(":8080", nil); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln("Server Failed", err)
	}
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()
	conn.WriteJSON("hi")
	conn.WriteJSON("there")
}
