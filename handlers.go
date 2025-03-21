package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// It should be a DB
var usernames = make(map[string]string)

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}
func saveUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	tmpl, err := template.ParseFiles("templates/chat.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}

	userId := uuid.New().String()
	usernames[userId] = username

	// Render the template with the message as data.
	err = tmpl.Execute(w, map[string]string{"userId": userId})
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
}
func wsHandler(room *ChatRoom, userId string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	log.Println("Client connected")

	username, ok := usernames[userId]
	if !ok {
		w.Write([]byte("userId not found"))
	}
	wsClient := NewWSClient(userId, username, conn, room)
	go wsClient.HandleWSConnection()
}
