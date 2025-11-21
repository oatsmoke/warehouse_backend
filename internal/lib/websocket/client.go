package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/oatsmoke/warehouse_backend/internal/lib/env"
	"github.com/oatsmoke/warehouse_backend/internal/lib/logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == env.GetClientUrl()
	},
}

type Client struct {
	conn     *websocket.Conn
	hub      *Hub
	location string
}

func NewClient(w http.ResponseWriter, r *http.Request, hub *Hub) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("failed websocket upgrade", err)
		return
	}

	client := &Client{
		conn: conn,
		hub:  hub,
	}

	client.hub.register <- client

	go func() {
		defer func() {
			client.hub.unregister <- client
			client.conn.Close()
		}()

		for {
			_, msg, err := client.conn.ReadMessage()
			if err != nil {
				logger.Error("failed read client message", err)
				break
			}
			client.location = string(msg)
		}
	}()
}
