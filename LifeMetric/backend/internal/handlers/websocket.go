package handlers

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Client represents an active clinician web dashboard connected over WebSockets
type Client struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}

// WebSocketHub manages active client connections and broadcasts incoming alerts
type WebSocketHub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	Mu         sync.RWMutex
}

// Global WebSocket Hub instance
var GlobalHub = &WebSocketHub{
	Clients:    make(map[*Client]bool),
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

// Run starts the WebSocket Hub's message routing loop (should run in a background goroutine)
func (h *WebSocketHub) Run() {
	log.Println("[WEBSOCKET HUB] Starting real-time message router...")
	for {
		select {
		case client := <-h.Register:
			h.Mu.Lock()
			h.Clients[client] = true
			h.Mu.Unlock()
			log.Printf("[WEBSOCKET HUB] Clinician connected. Active terminals: %d", len(h.Clients))

		case client := <-h.Unregister:
			h.Mu.Lock()
			if _, exists := h.Clients[client]; exists {
				delete(h.Clients, client)
				client.Conn.Close()
			}
			h.Mu.Unlock()
			log.Printf("[WEBSOCKET HUB] Clinician disconnected. Active terminals: %d", len(h.Clients))

		case message := <-h.Broadcast:
			h.Mu.RLock()
			for client := range h.Clients {
				go func(c *Client, msg []byte) {
					c.Mu.Lock()
					defer c.Mu.Unlock()
					if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
						log.Printf("[WEBSOCKET HUB] Error pushing alert to client: %v", err)
						h.Unregister <- c
					}
				}(client, message)
			}
			h.Mu.RUnlock()
		}
	}
}

// UpgradeWebSocketHandler upgrades incoming HTTP connections to WebSockets
func UpgradeWebSocketHandler(c *fiber.Ctx) error {
	// IsWebSocketUpgrade returns true if the client requested connection upgrade
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
