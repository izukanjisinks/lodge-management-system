package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"hr-system/internal/database"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type HolidayRepository struct {
	db *sql.DB
}

func NewHolidayRepository() *HolidayRepository {
	return &HolidayRepository{db: database.DB}
}

func (r *HolidayRepository) Create(h *models.Holiday) error {
	h.ID = uuid.New()
	now := time.Now()
	h.CreatedAt = now
	h.UpdatedAt = now
	_, err := r.db.Exec(`
		INSERT INTO holidays (id, name, date, description, is_recurring, location, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		h.ID, h.Name, h.Date, h.Description, h.IsRecurring, h.Location, h.IsActive, h.CreatedAt, h.UpdatedAt,
	)
	return err
}

func (r *HolidayRepository) GetByID(id uuid.UUID) (*models.Holiday, error) {
	return r.scanOne(r.db.QueryRow(`
		SELECT id, name, date, description, is_recurring, location, is_active, created_at, updated_at
		FROM holidays WHERE id=$1`, id))
}

func (r *HolidayRepository) List(year int, location string) ([]models.Holiday, error) {
	where := []string{}
	args := []interface{}{}
	i := 1

	if year > 0 {
		where = append(where, fmt.Sprintf("EXTRACT(YEAR FROM date)=$%d", i))
		args = append(args, year)
		i++
	}
	if location != "" {
		where = append(where, fmt.Sprintf("(location='' OR location=$%d)", i))
		args = append(args, location)
		i++
	}
	where = append(where, "is_active=TRUE")

	q := `SELECT id, name, date, description, is_recurring, location, is_active, created_at, updated_at FROM holidays`
	if len(where) > 0 {
		q += " WHERE " + strings.Join(where, " AND ")
	}
	q += " ORDER BY date"

	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Holiday
	for rows.Next() {
		h, err := r.scanOne(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *h)
	}
	return out, rows.Err()
}

func (r *HolidayRepository) Update(h *models.Holiday) error {
	h.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE holidays SET name=$1, date=$2, description=$3, is_recurring=$4, location=$5, is_active=$6, updated_at=$7
		WHERE id=$8`,
		h.Name, h.Date, h.Description, h.IsRecurring, h.Location, h.IsActive, h.UpdatedAt, h.ID,
	)
	return err
}

func (r *HolidayRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM holidays WHERE id=$1`, id)
	return err
}

// IsHoliday checks if a given date string (YYYY-MM-DD) is a holiday for the given location
func (r *HolidayRepository) IsHoliday(dateStr, location string) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(1) FROM holidays
		WHERE date=$1 AND is_active=TRUE AND (location='' OR location=$2)`,
		dateStr, location,
	).Scan(&count)
	return count > 0, err
}

// GetHolidaysInRange returns holiday dates in a given range (as time.Time slice)
func (r *HolidayRepository) GetHolidaysInRange(from, to time.Time, location string) (map[string]bool, error) {
	rows, err := r.db.Query(`
		SELECT date FROM holidays
		WHERE date>=$1 AND date<=$2 AND is_active=TRUE AND (location='' OR location=$3)`,
		from, to, location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	holidays := map[string]bool{}
	for rows.Next() {
		var d time.Time
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		holidays[d.Format("2006-01-02")] = true
	}
	return holidays, rows.Err()
}

func (r *HolidayRepository) scanOne(row rowScanner) (*models.Holiday, error) {
	var h models.Holiday
	err := row.Scan(&h.ID, &h.Name, &h.Date, &h.Description, &h.IsRecurring, &h.Location, &h.IsActive, &h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &h, nil
}
