package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/DPLakmal/guidenow-platform/go/internal/database"
	"github.com/DPLakmal/guidenow-platform/go/pkg/models"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for dev/demo
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client represents an active WebSocket connection
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	RoomID   string
	UserID   string
	UserName string
	Role     string
}

// Hub manages active WebSocket clients and broadcasts
type Hub struct {
	mu         sync.RWMutex
	rooms      map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *models.ChatMessage
	store      database.Store
}

func NewHub(store database.Store) *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *models.ChatMessage),
		store:      store,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, exists := h.rooms[client.RoomID]; !exists {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			h.mu.Unlock()

			// Broadcast user joined system message
			sysMsg := &models.ChatMessage{
				ID:         uuid.New().String(),
				RoomID:     client.RoomID,
				SenderID:   "system",
				SenderName: "GuideNow Bot",
				SenderRole: "system",
				Content:    client.UserName + " joined the chat.",
				Timestamp:  time.Now(),
			}
			h.BroadcastToRoom(client.RoomID, sysMsg)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.RoomID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.rooms, client.RoomID)
					}
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.BroadcastToRoom(msg.RoomID, msg)
		}
	}
}

func (h *Hub) BroadcastToRoom(roomID string, msg *models.ChatMessage) {
	h.mu.RLock()
	clients, ok := h.rooms[roomID]
	h.mu.RUnlock()

	if !ok {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	for client := range clients {
		select {
		case client.Send <- data:
		default:
			close(client.Send)
			delete(clients, client)
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(4096)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var incoming struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal(message, &incoming); err != nil {
			continue
		}

		if incoming.Content == "" {
			continue
		}

		chatMsg := &models.ChatMessage{
			ID:         uuid.New().String(),
			RoomID:     c.RoomID,
			SenderID:   c.UserID,
			SenderName: c.UserName,
			SenderRole: c.Role,
			Content:    incoming.Content,
			Timestamp:  time.Now(),
		}

		// Persist message
		if c.Hub.store != nil {
			_ = c.Hub.store.SaveChatMessage(chatMsg)
		}

		c.Hub.broadcast <- chatMsg
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles incoming websocket upgrade requests
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Upgrade error: %v", err)
		return
	}

	q := r.URL.Query()
	roomID := q.Get("room")
	if roomID == "" {
		roomID = "general"
	}
	userID := q.Get("userId")
	if userID == "" {
		userID = uuid.New().String()
	}
	userName := q.Get("name")
	if userName == "" {
		userName = "Traveler"
	}
	role := q.Get("role")
	if role == "" {
		role = "tourist"
	}

	client := &Client{
		Hub:      hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		RoomID:   roomID,
		UserID:   userID,
		UserName: userName,
		Role:     role,
	}

	client.Hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}
