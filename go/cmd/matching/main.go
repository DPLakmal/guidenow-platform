package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/DPLakmal/guidenow-platform/go/internal/config"
	"github.com/DPLakmal/guidenow-platform/go/internal/database"
	"github.com/DPLakmal/guidenow-platform/go/internal/redis"
	"github.com/DPLakmal/guidenow-platform/go/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()
	store := database.NewMemoryStore()
	redisClient := redis.NewClient(cfg.RedisURL)

	// Pre-seed Redis Geospatial index with guide locations
	ctx := context.Background()
	for _, g := range store.GetGuides() {
		_ = redisClient.AddLocation(ctx, g.ID, g.CurrentLat, g.CurrentLng)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "UP", "service": "matching-service"})
	})

	r.Route("/api/v1/matching", func(r chi.Router) {
		// GET /api/v1/matching/nearby?lat=6.0329&lng=80.2168&radius=10
		r.Get("/nearby", func(w http.ResponseWriter, r *http.Request) {
			latStr := r.URL.Query().Get("lat")
			lngStr := r.URL.Query().Get("lng")
			radiusStr := r.URL.Query().Get("radius")
			availableOnly := r.URL.Query().Get("available") != "false" // default true
			lang := r.URL.Query().Get("language")
			specialty := r.URL.Query().Get("specialty")

			if latStr == "" || lngStr == "" {
				// Default to Colombo Fort coordinates if not specified
				latStr = "6.9271"
				lngStr = "79.8612"
			}

			lat, err := strconv.ParseFloat(latStr, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Invalid latitude"})
				return
			}

			lng, err := strconv.ParseFloat(lngStr, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Invalid longitude"})
				return
			}

			radius := 25.0
			if radiusStr != "" {
				if rVal, err := strconv.ParseFloat(radiusStr, 64); err == nil && rVal > 0 {
					radius = rVal
				}
			}

			matches, err := redisClient.SearchNearby(r.Context(), lat, lng, radius)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: err.Error()})
				return
			}

			var result []models.Guide
			for _, m := range matches {
				g, err := store.GetGuideByID(m.GuideID)
				if err != nil {
					continue
				}
				if availableOnly && !g.IsAvailable {
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

				guideCopy := *g
				guideCopy.DistanceKm = m.DistanceKm
				result = append(result, guideCopy)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: true,
				Data: map[string]interface{}{
					"origin": map[string]float64{"lat": lat, "lng": lng},
					"radius_km": radius,
					"count": len(result),
					"guides": result,
				},
			})
		})

		// POST /api/v1/matching/location - update live GPS
		r.Post("/location", func(w http.ResponseWriter, r *http.Request) {
			var body models.LocationUpdate
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Invalid body"})
				return
			}

			if body.GuideID == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "guide_id required"})
				return
			}

			_ = store.UpdateGuideLocation(body.GuideID, body.Latitude, body.Longitude)
			_ = redisClient.AddLocation(r.Context(), body.GuideID, body.Latitude, body.Longitude)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: true,
				Message: "Guide location updated in Redis Geospatial index",
			})
		})
	})

	port := cfg.MatchingPort
	log.Printf("[Matching Service] Starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Matching service failed: %v", err)
	}
}
