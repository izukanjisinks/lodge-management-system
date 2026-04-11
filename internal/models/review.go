package models

import (
	"time"

	"github.com/google/uuid"
)

// RatingLabel tiers based on overall score
func RatingLabel(score float64) string {
	switch {
	case score >= 4.5:
		return "Exceptional"
	case score >= 4.0:
		return "Impressive"
	case score >= 3.5:
		return "Good"
	case score >= 3.0:
		return "Satisfactory"
	default:
		return "Needs Improvement"
	}
}

type Review struct {
	ID          uuid.UUID `json:"id"`
	BookingID   uuid.UUID `json:"booking_id"`
	GuestID     uuid.UUID `json:"guest_id"`
	Facilities  float64   `json:"facilities"`
	Cleanliness float64   `json:"cleanliness"`
	Services    float64   `json:"services"`
	Comfort     float64   `json:"comfort"`
	Location    float64   `json:"location"`
	Comment     string    `json:"comment,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type SubmitReviewRequest struct {
	BookingID   uuid.UUID `json:"booking_id"`
	Facilities  float64   `json:"facilities"`
	Cleanliness float64   `json:"cleanliness"`
	Services    float64   `json:"services"`
	Comfort     float64   `json:"comfort"`
	Location    float64   `json:"location"`
	Comment     string    `json:"comment,omitempty"`
}

type RatingCategory struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}

type RatingSummary struct {
	OverallScore float64          `json:"overall_score"`
	TotalReviews int              `json:"total_reviews"`
	Label        string           `json:"label"`
	Categories   []RatingCategory `json:"categories"`
}
