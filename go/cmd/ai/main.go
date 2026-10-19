package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DPLakmal/guidenow-platform/go/internal/ai"
	"github.com/DPLakmal/guidenow-platform/go/internal/config"
	"github.com/DPLakmal/guidenow-platform/go/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()
	parser := ai.NewParser(cfg.GeminiAPIKey)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "UP", "service": "ai-service"})
	})

	r.Route("/api/v1/ai", func(r chi.Router) {
		r.Post("/parse-query", func(w http.ResponseWriter, r *http.Request) {
			var body struct {
				Query string `json:"query"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Invalid JSON body"})
				return
			}

			filter, err := parser.Parse(r.Context(), body.Query)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: err.Error()})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: true,
				Data:    filter,
			})
		})

		r.Get("/examples", func(w http.ResponseWriter, r *http.Request) {
			examples := []string{
				"English or German speaking guide in Galle Fort for history and architecture under $25/hr",
				"Surfing instructor in Weligama under $35 available right now",
				"Temple of the Tooth and botanical garden cultural tour in Kandy",
				"French guide in Sigiriya for sunrise climb and rock art",
				"Tea plantation and Nine Arch Bridge mountain hiking guide in Ella",
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: true,
				Data:    examples,
			})
		})
	})

	port := cfg.AIPort
	log.Printf("[AI Service] Starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("AI service failed: %v", err)
	}
}
