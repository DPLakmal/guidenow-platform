package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/DPLakmal/guidenow-platform/go/internal/config"
	"github.com/DPLakmal/guidenow-platform/go/internal/database"
	"github.com/DPLakmal/guidenow-platform/go/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()
	store := database.NewMemoryStore()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "UP", "service": "guide-service"})
	})

	r.Route("/api/v1/guides", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			city := r.URL.Query().Get("city")
			lang := r.URL.Query().Get("language")
			specialty := r.URL.Query().Get("specialty")
			availableOnly := r.URL.Query().Get("available") == "true"
			maxRateStr := r.URL.Query().Get("max_rate")

			var maxRate float64
			if maxRateStr != "" {
				maxRate, _ = strconv.ParseFloat(maxRateStr, 64)
			}

			all := store.GetGuides()
			var filtered []models.Guide

			for _, g := range all {
				if availableOnly && !g.IsAvailable {
					continue
				}
				if city != "" && !strings.EqualFold(g.City, city) {
					continue
				}
				if maxRate > 0 && g.HourlyRate > maxRate {
					continue
				}
				if lang != "" {
					hasLang := false
					for _, l := range g.Languages {
						if strings.EqualFold(l, lang) {
							hasLang = true
							break
						}
					}
					if !hasLang {
						continue
					}
				}
				if specialty != "" {
					hasSpec := false
					for _, s := range g.Specialties {
						if strings.Contains(strings.ToLower(s), strings.ToLower(specialty)) {
							hasSpec = true
							break
						}
					}
					if !hasSpec {
						continue
					}
				}
				filtered = append(filtered, g)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: true,
				Data:    filtered,
			})
		})

		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			g, err := store.GetGuideByID(id)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(models.APIResponse{
					Success: false,
					Error:   err.Error(),
				})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: true,
				Data:    g,
			})
		})

		r.Patch("/{id}/availability", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			var body struct {
				Available bool `json:"available"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(models.APIResponse{
					Success: false,
					Error:   "Invalid request body",
				})
				return
			}

			updated, err := store.UpdateGuideAvailability(id, body.Available)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(models.APIResponse{
					Success: false,
					Error:   err.Error(),
				})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: true,
				Message: "Guide availability updated successfully",
				Data:    updated,
			})
		})
	})

	port := cfg.GuidePort
	log.Printf("[Guide Service] Starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Guide service failed: %v", err)
	}
}
