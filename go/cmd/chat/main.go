package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DPLakmal/guidenow-platform/go/internal/config"
	"github.com/DPLakmal/guidenow-platform/go/internal/database"
	"github.com/DPLakmal/guidenow-platform/go/internal/websocket"
	"github.com/DPLakmal/guidenow-platform/go/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

func main() {
	cfg := config.Load()
	store := database.NewMemoryStore()
	hub := websocket.NewHub(store)
	go hub.Run()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "UP", "service": "chat-service"})
	})

	// WebSocket handler
	r.Get("/ws/chat", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	r.Route("/api/v1/chat", func(r chi.Router) {
		r.Get("/rooms/{roomId}/messages", func(w http.ResponseWriter, r *http.Request) {
			roomID := chi.URLParam(r, "roomId")
			history := store.GetChatHistory(roomID)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: true,
				Data:    history,
			})
		})

		r.Post("/rooms/{roomId}/messages", func(w http.ResponseWriter, r *http.Request) {
			roomID := chi.URLParam(r, "roomId")
			var body struct {
				SenderID   string `json:"sender_id"`
				SenderName string `json:"sender_name"`
				SenderRole string `json:"sender_role"`
				Content    string `json:"content"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Invalid body"})
				return
			}

			msg := &models.ChatMessage{
				ID:         uuid.New().String(),
				RoomID:     roomID,
				SenderID:   body.SenderID,
				SenderName: body.SenderName,
				SenderRole: body.SenderRole,
				Content:    body.Content,
				Timestamp:  time.Now(),
			}

			_ = store.SaveChatMessage(msg)
			hub.BroadcastToRoom(roomID, msg)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: true,
				Data:    msg,
			})
		})
	})

	port := cfg.ChatPort
	log.Printf("[Chat Service] Starting on :%s (WebSocket: ws://localhost:%s/ws/chat)", port, port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Chat service failed: %v", err)
	}
}
