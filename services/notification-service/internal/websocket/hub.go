package websocket

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/notification-service/internal/models"
)

type Hub struct {
	clients    map[string]map[*websocket.Conn]bool // userID -> connections
	broadcast  chan *BroadcastMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	log        *logger.Logger
}

type Client struct {
	UserID string
	Conn   *websocket.Conn
}

type BroadcastMessage struct {
	UserID       string
	Notification *models.Notification
}

func NewHub(log *logger.Logger) *Hub {
	return &Hub{
		clients:    make(map[string]map[*websocket.Conn]bool),
		broadcast:  make(chan *BroadcastMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		log:        log,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*websocket.Conn]bool)
			}
			h.clients[client.UserID][client.Conn] = true
			h.mu.Unlock()
			h.log.Sugar().Infow("Client registered", "user_id", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if connections, ok := h.clients[client.UserID]; ok {
				if _, exists := connections[client.Conn]; exists {
					delete(connections, client.Conn)
					client.Conn.Close()
					if len(connections) == 0 {
						delete(h.clients, client.UserID)
					}
				}
			}
			h.mu.Unlock()
			h.log.Sugar().Infow("Client unregistered", "user_id", client.UserID)

		case message := <-h.broadcast:
			h.mu.RLock()
			connections := h.clients[message.UserID]
			h.mu.RUnlock()

			if connections != nil {
				data, err := json.Marshal(message.Notification)
				if err != nil {
					h.log.Sugar().Errorw("Failed to marshal notification", "error", err)
					continue
				}

				for conn := range connections {
					err := conn.WriteMessage(websocket.TextMessage, data)
					if err != nil {
						h.log.Sugar().Warnw("Failed to send message", "error", err)
						h.unregister <- &Client{UserID: message.UserID, Conn: conn}
					}
				}
			}
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(userID string, notification *models.Notification) {
	h.broadcast <- &BroadcastMessage{
		UserID:       userID,
		Notification: notification,
	}
}
