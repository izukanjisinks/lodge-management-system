package repository

import (
	"database/sql"
	"fmt"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type MealPlanRepository struct {
	db *sql.DB
}

func NewMealPlanRepository() *MealPlanRepository {
	return &MealPlanRepository{db: database.DB}
}

func (r *MealPlanRepository) Create(m *models.MealPlan) error {
	m.ID = uuid.New()
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	_, err := r.db.Exec(`
		INSERT INTO meal_plans (id, name, price_per_person_per_night, includes, description, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		m.ID, m.Name, m.PricePerPersonPerNight, pq.Array(m.Includes),
		m.Description, m.IsActive, m.CreatedAt, m.UpdatedAt,
	)
	return err
}

func (r *MealPlanRepository) GetByID(id uuid.UUID) (*models.MealPlan, error) {
	row := r.db.QueryRow(`
		SELECT id, name, price_per_person_per_night, includes, description, is_active, created_at, updated_at
		FROM meal_plans WHERE id = $1`, id)
	return scanMealPlan(row)
}

func (r *MealPlanRepository) List(isActive *bool, page, pageSize int) ([]models.MealPlan, int, error) {
	where := "1=1"
	args := []interface{}{}
	i := 1

	if isActive != nil {
		where = fmt.Sprintf("is_active = $%d", i)
		args = append(args, *isActive)
		i++
	}

	var total int
	if err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM meal_plans WHERE %s`, where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT id, name, price_per_person_per_night, includes, description, is_active, created_at, updated_at
		FROM meal_plans WHERE %s
		ORDER BY name ASC
		LIMIT $%d OFFSET $%d`, where, i, i+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var plans []models.MealPlan
	for rows.Next() {
		m, err := scanMealPlan(rows)
		if err != nil {
			return nil, 0, err
		}
		plans = append(plans, *m)
	}
	return plans, total, rows.Err()
}

func (r *MealPlanRepository) Update(m *models.MealPlan) error {
	m.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE meal_plans
		SET name=$1, price_per_person_per_night=$2, includes=$3, description=$4, is_active=$5, updated_at=$6
		WHERE id=$7`,
		m.Name, m.PricePerPersonPerNight, pq.Array(m.Includes),
		m.Description, m.IsActive, m.UpdatedAt, m.ID,
	)
	return err
}

func (r *MealPlanRepository) Delete(id uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM meal_plans WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("meal plan not found")
	}
	return nil
}

type mealPlanScanner interface {
	Scan(dest ...interface{}) error
}

func scanMealPlan(row mealPlanScanner) (*models.MealPlan, error) {
	var m models.MealPlan
	var description sql.NullString
	err := row.Scan(
		&m.ID, &m.Name, &m.PricePerPersonPerNight, pq.Array(&m.Includes),
		&description, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if description.Valid {
		m.Description = description.String
	}
	if m.Includes == nil {
		m.Includes = []string{}
	}
	return &m, nil
}
