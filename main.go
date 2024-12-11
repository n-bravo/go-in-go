package main

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/n-bravo/go-in-go/server"
)

func main() {
	webSocketHandler := server.WebSocketHandler{
		Upgrader: websocket.Upgrader{},
	}
	http.Handle("/", webSocketHandler)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
