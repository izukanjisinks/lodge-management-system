package repository

import (
	"database/sql"
	"fmt"
	"time"

	"hr-system/internal/database"
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type EmployeeDocumentRepository struct {
	db *sql.DB
}

func NewEmployeeDocumentRepository() *EmployeeDocumentRepository {
	return &EmployeeDocumentRepository{db: database.DB}
}

func (r *EmployeeDocumentRepository) Create(doc *models.EmployeeDocument) error {
	doc.ID = uuid.New()
	now := time.Now()
	doc.CreatedAt = now
	doc.UpdatedAt = now
	_, err := r.db.Exec(`
		INSERT INTO employee_documents
		(id, employee_id, document_type, title, description, file_url, file_name, file_size, mime_type,
		 uploaded_by, expiry_date, is_verified, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`,
		doc.ID, doc.EmployeeID, doc.DocumentType, doc.Title, doc.Description, doc.FileURL, doc.FileName,
		doc.FileSize, doc.MimeType, doc.UploadedBy, doc.ExpiryDate, doc.IsVerified, doc.CreatedAt, doc.UpdatedAt,
	)
	return err
}

func (r *EmployeeDocumentRepository) GetByID(id uuid.UUID) (*models.EmployeeDocument, error) {
	var d models.EmployeeDocument
	var verifiedBy sql.NullString
	err := r.db.QueryRow(`
		SELECT id, employee_id, document_type, title, description, file_url, file_name, file_size, mime_type,
		       uploaded_by, expiry_date, is_verified, verified_by, verified_at, created_at, updated_at, deleted_at
		FROM employee_documents WHERE id=$1 AND deleted_at IS NULL`, id,
	).Scan(&d.ID, &d.EmployeeID, &d.DocumentType, &d.Title, &d.Description, &d.FileURL, &d.FileName,
		&d.FileSize, &d.MimeType, &d.UploadedBy, &d.ExpiryDate, &d.IsVerified, &verifiedBy, &d.VerifiedAt,
		&d.CreatedAt, &d.UpdatedAt, &d.DeletedAt)
	if err != nil {
		return nil, err
	}
	if verifiedBy.Valid {
		p, _ := uuid.Parse(verifiedBy.String)
		d.VerifiedBy = &p
	}
	return &d, nil
}

func (r *EmployeeDocumentRepository) ListByEmployee(employeeID uuid.UUID, docType string) ([]models.EmployeeDocument, error) {
	query := `
		SELECT id, employee_id, document_type, title, description, file_url, file_name, file_size, mime_type,
		       uploaded_by, expiry_date, is_verified, verified_by, verified_at, created_at, updated_at, deleted_at
		FROM employee_documents WHERE employee_id=$1 AND deleted_at IS NULL`
	args := []interface{}{employeeID}
	if docType != "" {
		query += fmt.Sprintf(" AND document_type=$%d", 2)
		args = append(args, docType)
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []models.EmployeeDocument
	for rows.Next() {
		var d models.EmployeeDocument
		var verifiedBy sql.NullString
		if err := rows.Scan(&d.ID, &d.EmployeeID, &d.DocumentType, &d.Title, &d.Description, &d.FileURL, &d.FileName,
			&d.FileSize, &d.MimeType, &d.UploadedBy, &d.ExpiryDate, &d.IsVerified, &verifiedBy, &d.VerifiedAt,
			&d.CreatedAt, &d.UpdatedAt, &d.DeletedAt); err != nil {
			return nil, err
		}
		if verifiedBy.Valid {
			p, _ := uuid.Parse(verifiedBy.String)
			d.VerifiedBy = &p
		}
		docs = append(docs, d)
	}
	return docs, rows.Err()
}

func (r *EmployeeDocumentRepository) Verify(id uuid.UUID, verifiedBy uuid.UUID) error {
	now := time.Now()
	_, err := r.db.Exec(`
		UPDATE employee_documents SET is_verified=true, verified_by=$1, verified_at=$2, updated_at=$2
		WHERE id=$3 AND deleted_at IS NULL`,
		verifiedBy, now, id,
	)
	return err
}

func (r *EmployeeDocumentRepository) SoftDelete(id uuid.UUID) error {
	_, err := r.db.Exec(`UPDATE employee_documents SET deleted_at=$1, updated_at=$1 WHERE id=$2 AND deleted_at IS NULL`, time.Now(), id)
	return err
}
