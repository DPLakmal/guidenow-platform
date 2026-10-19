package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/DPLakmal/guidenow-platform/go/pkg/models"
	_ "github.com/lib/pq"
)

// Store defines database operations
type Store interface {
	GetGuides() []models.Guide
	GetGuideByID(id string) (*models.Guide, error)
	UpdateGuideAvailability(id string, available bool) (*models.Guide, error)
	UpdateGuideLocation(id string, lat, lng float64) error

	CreateBooking(req *models.TourRequest) error
	GetBookingByID(id string) (*models.TourRequest, error)
	UpdateBookingStatus(id string, status models.TourStatus) (*models.TourRequest, error)
	ListBookings() []models.TourRequest

	SaveChatMessage(msg *models.ChatMessage) error
	GetChatHistory(roomID string) []models.ChatMessage
}

// MemoryStore provides in-memory thread-safe storage initialized with seed data
type MemoryStore struct {
	mu       sync.RWMutex
	guides   map[string]models.Guide
	bookings map[string]models.TourRequest
	messages map[string][]models.ChatMessage
}

func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		guides:   make(map[string]models.Guide),
		bookings: make(map[string]models.TourRequest),
		messages: make(map[string][]models.ChatMessage),
	}

	// Seed with curated Sri Lanka guides
	for _, g := range DefaultSriLankaGuides() {
		store.guides[g.ID] = g
	}

	return store
}

func (s *MemoryStore) GetGuides() []models.Guide {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]models.Guide, 0, len(s.guides))
	for _, g := range s.guides {
		out = append(out, g)
	}
	return out
}

func (s *MemoryStore) GetGuideByID(id string) (*models.Guide, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	g, exists := s.guides[id]
	if !exists {
		return nil, fmt.Errorf("guide not found: %s", id)
	}
	return &g, nil
}

func (s *MemoryStore) UpdateGuideAvailability(id string, available bool) (*models.Guide, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	g, exists := s.guides[id]
	if !exists {
		return nil, fmt.Errorf("guide not found: %s", id)
	}
	g.IsAvailable = available
	g.UpdatedAt = time.Now()
	s.guides[id] = g
	return &g, nil
}

func (s *MemoryStore) UpdateGuideLocation(id string, lat, lng float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	g, exists := s.guides[id]
	if !exists {
		return fmt.Errorf("guide not found: %s", id)
	}
	g.CurrentLat = lat
	g.CurrentLng = lng
	g.UpdatedAt = time.Now()
	s.guides[id] = g
	return nil
}

func (s *MemoryStore) CreateBooking(req *models.TourRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	s.bookings[req.ID] = *req
	return nil
}

func (s *MemoryStore) GetBookingByID(id string) (*models.TourRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, exists := s.bookings[id]
	if !exists {
		return nil, fmt.Errorf("booking not found: %s", id)
	}
	return &b, nil
}

func (s *MemoryStore) UpdateBookingStatus(id string, status models.TourStatus) (*models.TourRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, exists := s.bookings[id]
	if !exists {
		return nil, fmt.Errorf("booking not found: %s", id)
	}
	b.Status = status
	b.UpdatedAt = time.Now()
	s.bookings[id] = b
	return &b, nil
}

func (s *MemoryStore) ListBookings() []models.TourRequest {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]models.TourRequest, 0, len(s.bookings))
	for _, b := range s.bookings {
		out = append(out, b)
	}
	return out
}

func (s *MemoryStore) SaveChatMessage(msg *models.ChatMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}
	s.messages[msg.RoomID] = append(s.messages[msg.RoomID], *msg)
	return nil
}

func (s *MemoryStore) GetChatHistory(roomID string) []models.ChatMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msgs, exists := s.messages[roomID]
	if !exists {
		return []models.ChatMessage{}
	}
	out := make([]models.ChatMessage, len(msgs))
	copy(out, msgs)
	return out
}

// InitPostgres attempts to connect to PostgreSQL, runs DDL, and returns a DB handle if successful
func InitPostgres(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS guides (
		id VARCHAR(64) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		bio TEXT,
		avatar TEXT,
		phone VARCHAR(50),
		email VARCHAR(100),
		city VARCHAR(100),
		hourly_rate NUMERIC(10, 2),
		rating NUMERIC(3, 2),
		review_count INT,
		languages TEXT[],
		specialties TEXT[],
		current_lat DOUBLE PRECISION,
		current_lng DOUBLE PRECISION,
		is_available BOOLEAN DEFAULT true,
		verified_status BOOLEAN DEFAULT true,
		license_number VARCHAR(100),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS bookings (
		id VARCHAR(64) PRIMARY KEY,
		tourist_id VARCHAR(64) NOT NULL,
		tourist_name VARCHAR(255),
		tourist_phone VARCHAR(50),
		guide_id VARCHAR(64) REFERENCES guides(id),
		guide_name VARCHAR(255),
		status VARCHAR(50) NOT NULL,
		meeting_point TEXT,
		destination TEXT,
		duration_hours NUMERIC(4, 2),
		hourly_rate NUMERIC(10, 2),
		total_estimated NUMERIC(10, 2),
		special_notes TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS chat_messages (
		id VARCHAR(64) PRIMARY KEY,
		room_id VARCHAR(64) NOT NULL,
		sender_id VARCHAR(64) NOT NULL,
		sender_name VARCHAR(255),
		sender_role VARCHAR(50),
		content TEXT NOT NULL,
		timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		read BOOLEAN DEFAULT false
	);
	`
	if _, err := db.Exec(schema); err != nil {
		log.Printf("[DB] Warning: Schema initialization query error: %v", err)
	}

	return db, nil
}
