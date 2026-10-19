package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DPLakmal/guidenow-platform/go/internal/config"
	"github.com/DPLakmal/guidenow-platform/go/internal/database"
	"github.com/DPLakmal/guidenow-platform/go/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

func main() {
	cfg := config.Load()
	store := database.NewMemoryStore()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "UP", "service": "booking-service"})
	})

	r.Route("/api/v1/bookings", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			list := store.ListBookings()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{Success: true, Data: list})
		})

		r.Post("/request", func(w http.ResponseWriter, r *http.Request) {
			var body struct {
				TouristID     string  `json:"tourist_id"`
				TouristName   string  `json:"tourist_name"`
				TouristPhone  string  `json:"tourist_phone"`
				GuideID       string  `json:"guide_id"`
				MeetingPoint  string  `json:"meeting_point"`
				Destination   string  `json:"destination"`
				DurationHours float64 `json:"duration_hours"`
				SpecialNotes  string  `json:"special_notes"`
			}

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Invalid JSON body"})
				return
			}

			if body.GuideID == "" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "guide_id is required"})
				return
			}

			guide, err := store.GetGuideByID(body.GuideID)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Guide not found"})
				return
			}

			duration := body.DurationHours
			if duration <= 0 {
				duration = 2.0 // default 2 hours quick tour
			}

			tour := &models.TourRequest{
				ID:             "tour-" + uuid.New().String()[:8],
				TouristID:      body.TouristID,
				TouristName:    body.TouristName,
				TouristPhone:   body.TouristPhone,
				GuideID:        body.GuideID,
				GuideName:      guide.Name,
				Status:         models.StatusRequested,
				MeetingPoint:   body.MeetingPoint,
				Destination:    body.Destination,
				DurationHours:  duration,
				HourlyRate:     guide.HourlyRate,
				TotalEstimated: duration * guide.HourlyRate,
				SpecialNotes:   body.SpecialNotes,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			if err := store.CreateBooking(tour); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: err.Error()})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: true,
				Message: "Tour request dispatched to guide",
				Data:    tour,
			})
		})

		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			booking, err := store.GetBookingByID(id)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Booking not found"})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{Success: true, Data: booking})
		})

		r.Post("/{id}/accept", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			updated, err := store.UpdateBookingStatus(id, models.StatusAccepted)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: true,
				Message: "Tour request accepted by guide",
				Data:    updated,
			})
		})

		r.Post("/{id}/start", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			updated, err := store.UpdateBookingStatus(id, models.StatusInProgress)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: true,
				Message: "Tour is now in progress",
				Data:    updated,
			})
		})

		r.Post("/{id}/complete", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			updated, err := store.UpdateBookingStatus(id, models.StatusCompleted)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: err.Error()})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.APIResponse{
				Success: true,
				Message: "Tour completed successfully",
				Data:    updated,
			})
		})
	})

	port := cfg.BookingPort
	log.Printf("[Booking Service] Starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Booking service failed: %v", err)
	}
}
