package websocket

import (
	"log"
	"sync"
)

type Hub struct {
	// registered clients per user ID
	clients map[string]map[*Client]bool
	
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; !ok {
				h.clients[client.userID] = make(map[*Client]bool)
			}
			h.clients[client.userID][client] = true
			h.mu.Unlock()
			log.Printf("[WS] User %s connected. Total clients for user: %d", client.userID, len(h.clients[client.userID]))

		case client := <-h.unregister:
			h.mu.Lock()
			if conns, ok := h.clients[client.userID]; ok {
				if _, ok := conns[client]; ok {
					delete(conns, client)
					close(client.send)
					if len(conns) == 0 {
						delete(h.clients, client.userID)
					}
					log.Printf("[WS] User %s disconnected.", client.userID)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) SendToUser(userID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if conns, ok := h.clients[userID]; ok {
		for client := range conns {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(conns, client)
			}
		}
	}
}
