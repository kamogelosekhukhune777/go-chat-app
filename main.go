package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	//only this orin is allowed
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:8080" || r.Host == "localhost:8080"
	},
}

func echo(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade: ", err)
		return
	}
	defer conn.Close()

	// Loop to handle incoming messages.
	for {
		// Read message from the client.
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("failed to read message: ", err)
			return
		}
		log.Printf("Received: %v", p)

		// Echo the message back to the client.
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			log.Println("failed to write message ", err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/echo", echo)
	log.Println("Server Starting on Port :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Failed to Start Server %v", err)
	}
}
