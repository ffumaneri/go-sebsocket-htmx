package main

import (
	"fmt"
	"log"
	"sync"
)

type ChatRoom struct {
	clients      map[string]*WSClient
	incomingChan chan Message
	senderChan   chan Message
	register     chan *WSClient
	unregister   chan *WSClient
	mutex        sync.RWMutex
}

func Room() *ChatRoom {
	return &ChatRoom{
		clients:      make(map[string]*WSClient),
		incomingChan: make(chan Message, 256),
		senderChan:   make(chan Message, 256),
		register:     make(chan *WSClient),
		unregister:   make(chan *WSClient),
	}
}
func (h *ChatRoom) Run() {
	log.Println("ChatRoom started")
	for {
		select {
		case client, ok := <-h.register:
			if !ok {
				return
			}
			h.mutex.Lock()
			h.clients[client.ID] = client
			h.mutex.Unlock()
			msg := Message{ClientID: client.ID, Username: "Admin", Text: fmt.Sprintf("%s has joined", client.username), IsAdmin: true}
			h.incomingChan <- msg
		case client, ok := <-h.unregister:
			if !ok {
				return
			}
			h.mutex.Lock()
			delete(h.clients, client.ID)
			h.mutex.Unlock()
		case msg, ok := <-h.incomingChan:
			if !ok {
				return
			}
			for _, client := range h.clients {
				client.room.senderChan <- msg
			}
		}
	}
}
