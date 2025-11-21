package websocket

import (
	"fmt"

	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
)

type Hub struct {
	clients    map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
			logger.Info(fmt.Sprintf("register new client %v", client.conn.RemoteAddr()))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				logger.Info(fmt.Sprintf("unregister client %v", client.conn.RemoteAddr()))
			}
		}
	}
}
