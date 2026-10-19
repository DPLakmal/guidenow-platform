package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/DPLakmal/guidenow-platform/go/internal/config"
	"github.com/DPLakmal/guidenow-platform/go/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func newReverseProxy(target string) http.Handler {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid target URL %s: %v", target, err)
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	return proxy
}

func main() {
	cfg := config.Load()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "UP",
			"gateway":  "guidenow-api-gateway",
			"services": map[string]string{
				"guide":    cfg.GuideServiceURL,
				"matching": cfg.MatchingServiceURL,
				"chat":     cfg.ChatServiceURL,
				"ai":       cfg.AIServiceURL,
				"booking":  cfg.BookingServiceURL,
			},
		})
	})

	r.Get("/api/v1/system/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.APIResponse{
			Success: true,
			Message: "GuideNow API Gateway online",
			Data: map[string]interface{}{
				"version": "1.0.0",
				"region":  "ap-south-1",
				"environment": "production",
			},
		})
	})

	// Reverse proxy routes
	guideProxy := newReverseProxy(cfg.GuideServiceURL)
	matchingProxy := newReverseProxy(cfg.MatchingServiceURL)
	chatProxy := newReverseProxy(cfg.ChatServiceURL)
	aiProxy := newReverseProxy(cfg.AIServiceURL)
	bookingProxy := newReverseProxy(cfg.BookingServiceURL)

	r.HandleFunc("/api/v1/guides*", func(w http.ResponseWriter, r *http.Request) {
		guideProxy.ServeHTTP(w, r)
	})

	r.HandleFunc("/api/v1/matching*", func(w http.ResponseWriter, r *http.Request) {
		matchingProxy.ServeHTTP(w, r)
	})

	r.HandleFunc("/api/v1/chat*", func(w http.ResponseWriter, r *http.Request) {
		chatProxy.ServeHTTP(w, r)
	})

	r.HandleFunc("/ws/*", func(w http.ResponseWriter, r *http.Request) {
		// WebSocket upgrade through proxy
		chatProxy.ServeHTTP(w, r)
	})

	r.HandleFunc("/api/v1/ai*", func(w http.ResponseWriter, r *http.Request) {
		aiProxy.ServeHTTP(w, r)
	})

	r.HandleFunc("/api/v1/bookings*", func(w http.ResponseWriter, r *http.Request) {
		bookingProxy.ServeHTTP(w, r)
	})

	port := cfg.GatewayPort
	log.Printf("[Gateway Service] Reverse Proxy listening on :%s", port)
	log.Printf("  -> /api/v1/guides/*   ➔ %s", cfg.GuideServiceURL)
	log.Printf("  -> /api/v1/matching/* ➔ %s", cfg.MatchingServiceURL)
	log.Printf("  -> /api/v1/chat/*     ➔ %s", cfg.ChatServiceURL)
	log.Printf("  -> /ws/*              ➔ %s", cfg.ChatServiceURL)
	log.Printf("  -> /api/v1/ai/*       ➔ %s", cfg.AIServiceURL)
	log.Printf("  -> /api/v1/bookings/* ➔ %s", cfg.BookingServiceURL)

	if err := http.ListenAndServe(":"+port, r); err != nil && !strings.Contains(err.Error(), "Server closed") {
		log.Fatalf("Gateway service failed: %v", err)
	}
}
