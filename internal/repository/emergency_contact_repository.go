package repository

import (
	"database/sql"
	"time"

	"hr-system/internal/database"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type EmergencyContactRepository struct {
	db *sql.DB
}

func NewEmergencyContactRepository() *EmergencyContactRepository {
	return &EmergencyContactRepository{db: database.DB}
}

func (r *EmergencyContactRepository) Create(c *models.EmergencyContact) error {
	c.ID = uuid.New()
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	_, err := r.db.Exec(`
		INSERT INTO emergency_contacts (id, employee_id, name, relationship, phone, email, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		c.ID, c.EmployeeID, c.Name, c.Relationship, c.Phone, c.Email, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *EmergencyContactRepository) GetByID(id uuid.UUID) (*models.EmergencyContact, error) {
	var c models.EmergencyContact
	err := r.db.QueryRow(`
		SELECT id, employee_id, name, relationship, phone, email, created_at, updated_at
		FROM emergency_contacts WHERE id=$1`, id,
	).Scan(&c.ID, &c.EmployeeID, &c.Name, &c.Relationship, &c.Phone, &c.Email, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *EmergencyContactRepository) ListByEmployee(employeeID uuid.UUID) ([]models.EmergencyContact, error) {
	rows, err := r.db.Query(`
		SELECT id, employee_id, name, relationship, phone, email, created_at, updated_at
		FROM emergency_contacts WHERE employee_id=$1 ORDER BY created_at`, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []models.EmergencyContact
	for rows.Next() {
		var c models.EmergencyContact
		if err := rows.Scan(&c.ID, &c.EmployeeID, &c.Name, &c.Relationship, &c.Phone, &c.Email, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, rows.Err()
}

func (r *EmergencyContactRepository) Update(c *models.EmergencyContact) error {
	c.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE emergency_contacts SET name=$1, relationship=$2, phone=$3, email=$4, updated_at=$5 WHERE id=$6`,
		c.Name, c.Relationship, c.Phone, c.Email, c.UpdatedAt, c.ID,
	)
	return err
}

func (r *EmergencyContactRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM emergency_contacts WHERE id=$1`, id)
	return err
}

// dummy to satisfy compiler when db is unused
var _ = (*sql.DB)(nil)
