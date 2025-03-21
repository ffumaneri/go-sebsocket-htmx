package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"time"
)

type Message struct {
	ClientID string `json:"clientID"`
	Text     string `json:"text"`
}

type WSClient struct {
	ID     string
	wsconn *websocket.Conn
	room   *ChatRoom
}

func NewWSClient(id string, wsconn *websocket.Conn, hub *ChatRoom) *WSClient {
	return &WSClient{ID: id, wsconn: wsconn, room: hub}
}

func (c *WSClient) HandleWSConnection() {
	c.room.register <- c

	go c.WSReader()
	go c.WSWriter()
	return
}

func (c *WSClient) WSReader() {
	defer func() {
		c.wsconn.Close()
		log.Printf("WSClient %s Disconnected\n", c.ID)
	}()

	c.wsconn.SetReadLimit(512)
	c.wsconn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.wsconn.SetPongHandler(func(string) error {
		c.wsconn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	log.Printf("WSClient %s Connected\n", c.ID)
	for {
		// Read message
		_, wsmsg, err := c.wsconn.ReadMessage()
		if err != nil {
			c.room.unregister <- c
		}
		message := &Message{}
		message.ClientID = c.ID
		reader := bytes.NewReader(wsmsg)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(&message)
		if err != nil {
			log.Println("Decode error:", err)
			break
		}
		c.room.incomingChan <- *message
	}
	return
}

func (c *WSClient) WSWriter() {
	ticker := time.NewTicker(59 * time.Second)
	defer func() {
		ticker.Stop()
		c.wsconn.Close()
		log.Printf("WSClient %s Disconnected\n", c.ID)
	}()
	for {
		select {
		case msg, ok := <-c.room.senderChan:
			c.wsconn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The room closed the channel.
				c.wsconn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.wsconn.WriteMessage(websocket.TextMessage, getMessageTemplate(&msg))
			if err != nil {
				return
			}
		case <-ticker.C:
			c.wsconn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.wsconn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func getMessageTemplate(msg *Message) []byte {
	tmpl, err := template.ParseFiles("templates/message.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}

	// Render the template with the message as data.
	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, msg)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}

	return renderedMessage.Bytes()
}
