package tests

import (
	"context"
	"math"
	"testing"

	"github.com/DPLakmal/guidenow-platform/go/internal/database"
	"github.com/DPLakmal/guidenow-platform/go/internal/redis"
)

func TestHaversineDistance(t *testing.T) {
	// Colombo Fort (6.9271, 79.8612) to Galle Fort (6.0329, 80.2168) is ~115-120 km
	colomboLat, colomboLng := 6.9271, 79.8612
	galleLat, galleLng := 6.0329, 80.2168

	dist := redis.Haversine(colomboLat, colomboLng, galleLat, galleLng)
	if dist < 100 || dist > 115 {
		t.Errorf("Expected straight-line distance between Colombo and Galle to be ~107km, got %.2f km", dist)
	}

	// Distance from point to itself must be 0
	zeroDist := redis.Haversine(colomboLat, colomboLng, colomboLat, colomboLng)
	if math.Abs(zeroDist) > 0.001 {
		t.Errorf("Expected 0 distance for same coordinates, got %.4f", zeroDist)
	}
}

func TestRedisGeospatialNearbySearch(t *testing.T) {
	client := redis.NewClient("localhost:9999") // Triggers fast in-memory fallback
	ctx := context.Background()

	// Seed guides
	guides := database.DefaultSriLankaGuides()
	for _, g := range guides {
		err := client.AddLocation(ctx, g.ID, g.CurrentLat, g.CurrentLng)
		if err != nil {
			t.Fatalf("Failed to add guide location: %v", err)
		}
	}

	// Search around Galle Fort (6.0329, 80.2168) with 10km radius
	matches, err := client.SearchNearby(ctx, 6.0329, 80.2168, 10.0)
	if err != nil {
		t.Fatalf("SearchNearby failed: %v", err)
	}

	if len(matches) < 2 {
		t.Fatalf("Expected at least 2 Galle guides within 10km, got %d", len(matches))
	}

	// Verify sorting order (nearest first)
	for i := 1; i < len(matches); i++ {
		if matches[i].DistanceKm < matches[i-1].DistanceKm {
			t.Errorf("Results not sorted ascending by distance: match[%d]=%.2f < match[%d]=%.2f",
				i, matches[i].DistanceKm, i-1, matches[i-1].DistanceKm)
		}
	}
}
