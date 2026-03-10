package repositories

import (
	"database/sql"

	"hr-system/internal/database"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type PasswordHistoryRepository struct {
	db *sql.DB
}

func NewPasswordHistoryRepository() *PasswordHistoryRepository {
	return &PasswordHistoryRepository{db: database.DB}
}

// Create adds a new password hash to the history
func (r *PasswordHistoryRepository) Create(history *models.PasswordHistory) error {
	query := `
		INSERT INTO password_history (id, user_id, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, history.ID, history.UserID, history.PasswordHash, history.CreatedAt)
	return err
}

// GetRecentByUserID retrieves the most recent password hashes for a user
// Limited to the last 10 passwords to prevent reuse
func (r *PasswordHistoryRepository) GetRecentByUserID(userID uuid.UUID, limit int) ([]models.PasswordHistory, error) {
	if limit <= 0 {
		limit = 10 // Default to last 10 passwords
	}

	query := `
		SELECT id, user_id, password_hash, created_at
		FROM password_history
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.PasswordHistory
	for rows.Next() {
		var h models.PasswordHistory
		if err := rows.Scan(&h.ID, &h.UserID, &h.PasswordHash, &h.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	return history, rows.Err()
}

// DeleteOldHistory removes password history entries older than the specified limit
// Keeps only the most recent N entries per user
func (r *PasswordHistoryRepository) DeleteOldHistory(userID uuid.UUID, keepCount int) error {
	query := `
		DELETE FROM password_history
		WHERE user_id = $1
		AND id NOT IN (
			SELECT id FROM password_history
			WHERE user_id = $1
			ORDER BY created_at DESC
			LIMIT $2
		)
	`
	_, err := r.db.Exec(query, userID, keepCount)
	return err
}
