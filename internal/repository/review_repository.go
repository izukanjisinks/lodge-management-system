package repository

import (
	"database/sql"
	"fmt"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type ReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{db: database.DB}
}

func (r *ReviewRepository) Create(review *models.Review) error {
	query := `
		INSERT INTO reviews (id, booking_id, guest_id, facilities, cleanliness, services, comfort, location, comment, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	review.ID = uuid.New()
	review.CreatedAt = time.Now()

	_, err := r.db.Exec(query,
		review.ID, review.BookingID, review.GuestID,
		review.Facilities, review.Cleanliness, review.Services, review.Comfort, review.Location,
		review.Comment, review.CreatedAt,
	)
	return err
}

func (r *ReviewRepository) ExistsByBookingID(bookingID uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) FROM reviews WHERE booking_id = $1`

	var count int
	err := r.db.QueryRow(query, bookingID).Scan(&count)
	return count > 0, err
}

func (r *ReviewRepository) GetSummary() (*models.RatingSummary, error) {
	query := `
		SELECT
		    COUNT(*)                                                              AS total_reviews,
		    ROUND(AVG((facilities + cleanliness + services + comfort + location) / 5.0)::numeric, 1) AS overall_score,
		    ROUND(AVG(facilities)::numeric,  1) AS avg_facilities,
		    ROUND(AVG(cleanliness)::numeric, 1) AS avg_cleanliness,
		    ROUND(AVG(services)::numeric,    1) AS avg_services,
		    ROUND(AVG(comfort)::numeric,     1) AS avg_comfort,
		    ROUND(AVG(location)::numeric,    1) AS avg_location
		FROM reviews`

	var total int
	var overall, facilities, cleanliness, services, comfort, location float64

	err := r.db.QueryRow(query).Scan(
		&total, &overall,
		&facilities, &cleanliness, &services, &comfort, &location,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to compute rating summary: %w", err)
	}

	summary := &models.RatingSummary{
		OverallScore: overall,
		TotalReviews: total,
		Label:        models.RatingLabel(overall),
		Categories: []models.RatingCategory{
			{Label: "Facilities", Score: facilities},
			{Label: "Cleanliness", Score: cleanliness},
			{Label: "Services", Score: services},
			{Label: "Comfort", Score: comfort},
			{Label: "Location", Score: location},
		},
	}
	return summary, nil
}
