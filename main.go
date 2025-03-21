package main

import (
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}
func saveUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	tmpl, err := template.ParseFiles("templates/chat.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}

	// Render the template with the message as data.
	err = tmpl.Execute(w, map[string]string{"username": username})
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
}
func wsHandler(hub *ChatRoom, userId string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	log.Println("Client connected")

	wsClient := NewWSClient(userId, conn, hub)
	go wsClient.HandleWSConnection()
}

func main() {
	hub := NewHub()
	go hub.Run()
	http.HandleFunc("/", homePage)
	http.HandleFunc("POST /save-username", saveUsername)
	http.HandleFunc("/ws/", func(writer http.ResponseWriter, r *http.Request) {
		path := r.URL.Path                // Get the entire path
		parts := strings.Split(path, "/") // Split the path into parts
		if len(parts) >= 3 {
			wsHandler(hub, parts[2], writer, r)
		}
	})
	log.Println("Listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
