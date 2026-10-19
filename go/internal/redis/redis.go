package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/DPLakmal/guidenow-platform/go/pkg/models"
	rdb "github.com/redis/go-redis/v9"
)

const (
	GeoKeyGuides = "guides:locations"
)

// GeoMatch represents a matched guide with distance
type GeoMatch struct {
	GuideID    string
	DistanceKm float64
	Latitude   float64
	Longitude  float64
}

// Client abstracts Redis Geospatial and Pub/Sub operations
type Client struct {
	rdb          *rdb.Client
	useFallback  bool
	fallbackMu   sync.RWMutex
	fallbackLocs map[string][2]float64 // guide_id -> [lat, lng]
}

func NewClient(addr string) *Client {
	client := &Client{
		fallbackLocs: make(map[string][2]float64),
	}

	opts := &rdb.Options{
		Addr:         addr,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	redisClient := rdb.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		// Switch to high-performance in-memory fallback
		client.useFallback = true
	} else {
		client.rdb = redisClient
	}

	return client
}

// AddLocation records or updates a guide's GPS coordinates
func (c *Client) AddLocation(ctx context.Context, guideID string, lat, lng float64) error {
	if c.useFallback || c.rdb == nil {
		c.fallbackMu.Lock()
		defer c.fallbackMu.Unlock()
		c.fallbackLocs[guideID] = [2]float64{lat, lng}
		return nil
	}

	return c.rdb.GeoAdd(ctx, GeoKeyGuides, &rdb.GeoLocation{
		Name:      guideID,
		Latitude:  lat,
		Longitude: lng,
	}).Err()
}

// SearchNearby queries guides within a radius (in kilometers) of a point
func (c *Client) SearchNearby(ctx context.Context, lat, lng float64, radiusKm float64) ([]GeoMatch, error) {
	if radiusKm <= 0 {
		radiusKm = 25.0 // default 25km radius
	}

	if c.useFallback || c.rdb == nil {
		c.fallbackMu.RLock()
		defer c.fallbackMu.RUnlock()

		var results []GeoMatch
		for id, coords := range c.fallbackLocs {
			gLat, gLng := coords[0], coords[1]
			dist := Haversine(lat, lng, gLat, gLng)
			if dist <= radiusKm {
				results = append(results, GeoMatch{
					GuideID:    id,
					DistanceKm: math.Round(dist*100) / 100,
					Latitude:   gLat,
					Longitude:  gLng,
				})
			}
		}

		sort.Slice(results, func(i, j int) bool {
			return results[i].DistanceKm < results[j].DistanceKm
		})
		return results, nil
	}

	locations, err := c.rdb.GeoSearchLocation(ctx, GeoKeyGuides, &rdb.GeoSearchLocationQuery{
		GeoSearchQuery: rdb.GeoSearchQuery{
			Longitude:  lng,
			Latitude:   lat,
			Radius:     radiusKm,
			RadiusUnit: "km",
			Sort:       "ASC",
		},
		WithCoord: true,
		WithDist:  true,
	}).Result()

	if err != nil {
		return nil, err
	}

	matches := make([]GeoMatch, 0, len(locations))
	for _, loc := range locations {
		matches = append(matches, GeoMatch{
			GuideID:    loc.Name,
			DistanceKm: math.Round(loc.Dist*100) / 100,
			Latitude:   loc.Latitude,
			Longitude:  loc.Longitude,
		})
	}

	return matches, nil
}

// PublishChatMessage broadcasts a message across instances
func (c *Client) PublishChatMessage(ctx context.Context, msg *models.ChatMessage) error {
	if c.useFallback || c.rdb == nil {
		return nil
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	channel := fmt.Sprintf("chat:room:%s", msg.RoomID)
	return c.rdb.Publish(ctx, channel, payload).Err()
}

// Haversine calculates distance between two lat/lng pairs in kilometers
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	radLat1 := lat1 * (math.Pi / 180.0)
	radLat2 := lat2 * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(radLat1)*math.Cos(radLat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}
