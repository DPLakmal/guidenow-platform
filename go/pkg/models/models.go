package models

import "time"

// Guide represents a registered local guide
type Guide struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Bio            string    `json:"bio"`
	Avatar         string    `json:"avatar"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	City           string    `json:"city"`
	HourlyRate     float64   `json:"hourly_rate"`
	Rating         float64   `json:"rating"`
	ReviewCount    int       `json:"review_count"`
	Languages      []string  `json:"languages"`
	Specialties    []string  `json:"specialties"`
	CurrentLat     float64   `json:"current_lat"`
	CurrentLng     float64   `json:"current_lng"`
	DistanceKm     float64   `json:"distance_km,omitempty"`
	IsAvailable    bool      `json:"is_available"`
	VerifiedStatus bool      `json:"verified_status"`
	LicenseNumber  string    `json:"license_number"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// GuideReview represents a tourist review for a guide
type GuideReview struct {
	ID        string    `json:"id"`
	GuideID   string    `json:"guide_id"`
	Tourist   string    `json:"tourist_name"`
	Rating    float64   `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// TourStatus describes the state of a tour booking
type TourStatus string

const (
	StatusRequested  TourStatus = "REQUESTED"
	StatusAccepted   TourStatus = "ACCEPTED"
	StatusDeclined   TourStatus = "DECLINED"
	StatusInProgress TourStatus = "IN_PROGRESS"
	StatusCompleted  TourStatus = "COMPLETED"
	StatusCancelled  TourStatus = "CANCELLED"
)

// TourRequest represents a booking request between a tourist and a guide
type TourRequest struct {
	ID             string     `json:"id"`
	TouristID      string     `json:"tourist_id"`
	TouristName    string     `json:"tourist_name"`
	TouristPhone   string     `json:"tourist_phone"`
	GuideID        string     `json:"guide_id"`
	GuideName      string     `json:"guide_name"`
	Status         TourStatus `json:"status"`
	MeetingPoint   string     `json:"meeting_point"`
	Destination    string     `json:"destination"`
	DurationHours  float64    `json:"duration_hours"`
	HourlyRate     float64    `json:"hourly_rate"`
	TotalEstimated float64    `json:"total_estimated"`
	SpecialNotes   string     `json:"special_notes,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ChatMessage represents a real-time message between tourist and guide
type ChatMessage struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	SenderID  string    `json:"sender_id"`
	SenderName string   `json:"sender_name"`
	SenderRole string   `json:"sender_role"` // "tourist" or "guide"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Read      bool      `json:"read"`
}

// LocationUpdate represents a live GPS ping from a guide or tourist
type LocationUpdate struct {
	GuideID   string    `json:"guide_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}

// AISearchFilter represents the structured output of parsing natural language
type AISearchFilter struct {
	OriginalQuery string   `json:"original_query"`
	Location      string   `json:"location,omitempty"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	RadiusKm      float64  `json:"radius_km,omitempty"`
	Languages     []string `json:"languages,omitempty"`
	Specialties   []string `json:"specialties,omitempty"`
	MaxRate       float64  `json:"max_rate,omitempty"`
	MinRating     float64  `json:"min_rating,omitempty"`
	ImmediateOnly bool     `json:"immediate_only"`
	Confidence    float64  `json:"confidence"`
	Explanation   string   `json:"explanation,omitempty"`
}

// APIResponse standardizes JSON responses
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
