package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	room := Room()
	http.HandleFunc("/", homePage)
	http.HandleFunc("POST /save-username", saveUsername)
	http.HandleFunc("/ws/", func(writer http.ResponseWriter, r *http.Request) {
		path := r.URL.Path                // Get the entire path
		parts := strings.Split(path, "/") // Split the path into parts
		if len(parts) >= 3 {
			wsHandler(room, parts[2], writer, r)
		}
	})

	go room.Run()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
	log.Println("Listening on port 8080")
}
