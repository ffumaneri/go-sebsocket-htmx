package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}
func wsHandler(hub *ChatRoom, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	log.Println("Client connected")

	wsClient := NewWSClient(uuid.New().String(), conn, hub)
	go wsClient.HandleWSConnection()
}

func main() {
	hub := NewHub()
	go hub.Run()
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		wsHandler(hub, writer, request)
	})
	log.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
