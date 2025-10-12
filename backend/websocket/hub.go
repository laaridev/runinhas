package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// MessageType defines the type of WebSocket message
type MessageType string

const (
	MessageTypeConfig   MessageType = "config"
	MessageTypeEvent    MessageType = "event"
	MessageTypeStatus   MessageType = "status"
	MessageTypePing     MessageType = "ping"
	MessageTypePong     MessageType = "pong"
)

// Message represents a WebSocket message
type Message struct {
	Type    MessageType    `json:"type"`
	Payload interface{}    `json:"payload"`
}

// Client represents a WebSocket client
type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
	Hub  *Hub
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client %s connected", client.ID)
			
			// Send initial status
			h.SendToClient(client, Message{
				Type:    MessageTypeStatus,
				Payload: map[string]interface{}{"connected": true},
			})

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				h.mu.Unlock()
				log.Printf("Client %s disconnected", client.ID)
			} else {
				h.mu.Unlock()
			}

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					// Client's send channel is full, close it
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all clients
func (h *Hub) Broadcast(msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	
	select {
	case h.broadcast <- data:
	default:
		// Broadcast channel is full
		log.Println("Broadcast channel is full")
	}
	
	return nil
}

// SendToClient sends a message to a specific client
func (h *Hub) SendToClient(client *Client, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	
	select {
	case client.Send <- data:
	default:
		// Client's send channel is full
		log.Printf("Client %s send channel is full", client.ID)
	}
	
	return nil
}

// ClientCount returns the number of connected clients
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// ReadPump pumps messages from websocket to hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	
	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}
		
		// Handle ping
		if msg.Type == MessageTypePing {
			c.Hub.SendToClient(c, Message{
				Type: MessageTypePong,
			})
			continue
		}
		
		// Process other messages here
		log.Printf("Received message from %s: %+v", c.ID, msg)
	}
}

// WritePump pumps messages from hub to websocket
func (c *Client) WritePump() {
	defer c.Conn.Close()
	
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
