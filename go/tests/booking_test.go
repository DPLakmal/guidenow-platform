package tests

import (
	"testing"
	"time"

	"github.com/DPLakmal/guidenow-platform/go/internal/database"
	"github.com/DPLakmal/guidenow-platform/go/pkg/models"
)

func TestBookingLifecycleTransitions(t *testing.T) {
	store := database.NewMemoryStore()

	// 1. Create a tour request
	req := &models.TourRequest{
		ID:             "tour-test-01",
		TouristID:      "user-100",
		TouristName:    "Sarah Connor",
		GuideID:        "guide-galle-01",
		GuideName:      "Chaminda Jayasuriya",
		Status:         models.StatusRequested,
		MeetingPoint:   "Galle Dutch Gate",
		Destination:    "Galle Fort Ramparts",
		DurationHours:  3.0,
		HourlyRate:     25.0,
		TotalEstimated: 75.0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := store.CreateBooking(req)
	if err != nil {
		t.Fatalf("Failed to create booking: %v", err)
	}

	// 2. Guide accepts
	accepted, err := store.UpdateBookingStatus(req.ID, models.StatusAccepted)
	if err != nil {
		t.Fatalf("Failed to accept booking: %v", err)
	}
	if accepted.Status != models.StatusAccepted {
		t.Errorf("Expected status %s, got %s", models.StatusAccepted, accepted.Status)
	}

	// 3. Tour starts
	inProgress, err := store.UpdateBookingStatus(req.ID, models.StatusInProgress)
	if err != nil {
		t.Fatalf("Failed to start booking: %v", err)
	}
	if inProgress.Status != models.StatusInProgress {
		t.Errorf("Expected status %s, got %s", models.StatusInProgress, inProgress.Status)
	}

	// 4. Tour completes
	completed, err := store.UpdateBookingStatus(req.ID, models.StatusCompleted)
	if err != nil {
		t.Fatalf("Failed to complete booking: %v", err)
	}
	if completed.Status != models.StatusCompleted {
		t.Errorf("Expected status %s, got %s", models.StatusCompleted, completed.Status)
	}
}
