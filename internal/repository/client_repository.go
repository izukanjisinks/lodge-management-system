package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository() *ClientRepository {
	return &ClientRepository{db: database.DB}
}

// ─── Individual ───────────────────────────────────────────────────────────────

func (r *ClientRepository) CreateIndividual(c *models.IndividualClient, orgID uuid.UUID) error {
	query := `
		INSERT INTO individual_profiles
		    (id, full_name, email, phone, id_passport_number, nationality, status, notes, org_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

	c.ID = uuid.New()
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	_, err := r.db.Exec(query,
		c.ID, c.FullName, c.Email, c.Phone, c.IDPassportNumber,
		c.Nationality, c.Status, c.Notes, orgID, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

// CreateIndividualTx inserts an individual profile within an existing transaction.
func (r *ClientRepository) CreateIndividualTx(tx *sql.Tx, c *models.IndividualClient, userID uuid.UUID) error {
	query := `
		INSERT INTO individual_profiles
		    (id, user_id, full_name, email, phone, id_passport_number, nationality, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,'active',$8,$9)`

	c.ID = uuid.New()
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	_, err := tx.Exec(query,
		c.ID, userID, c.FullName, c.Email, c.Phone,
		c.IDPassportNumber, c.Nationality, now, now,
	)
	return err
}

// CreateIndividualInTx inserts a new individual profile scoped to an org within an existing transaction.
func (r *ClientRepository) CreateIndividualInTx(tx *sql.Tx, c *models.IndividualClient, orgID uuid.UUID) error {
	c.ID = uuid.New()
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	_, err := tx.Exec(`
		INSERT INTO individual_profiles
		    (id, full_name, email, phone, id_passport_number, nationality, status, notes, org_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		c.ID, c.FullName, c.Email, c.Phone, c.IDPassportNumber,
		c.Nationality, c.Status, c.Notes, orgID, now, now,
	)
	return err
}

// CreateCorporateInTx inserts a new corporate profile scoped to an org within an existing transaction.
func (r *ClientRepository) CreateCorporateInTx(tx *sql.Tx, c *models.CorporateClient, orgID uuid.UUID) error {
	c.ID = uuid.New()
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	_, err := tx.Exec(`
		INSERT INTO corporate_profiles
		    (id, company_name, contact_person, email, phone, company_reg_number, industry, status, notes, org_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		c.ID, c.CompanyName, c.ContactPerson, c.Email, c.Phone,
		c.CompanyRegNumber, c.Industry, c.Status, c.Notes, orgID, now, now,
	)
	return err
}

// GetIndividualByUserID returns the individual profile linked to a user account.
func (r *ClientRepository) GetIndividualByUserID(userID uuid.UUID) (*models.IndividualClient, error) {
	query := `
		SELECT id, full_name, email, phone, id_passport_number, nationality, status, notes, created_at, updated_at
		FROM individual_profiles
		WHERE user_id = $1`
	return r.scanIndividual(r.db.QueryRow(query, userID))
}

// UpdateIndividualByUserID updates the profile fields a guest is allowed to change.
func (r *ClientRepository) UpdateIndividualByUserID(userID uuid.UUID, c *models.IndividualClient) error {
	query := `
		UPDATE individual_profiles
		SET full_name=$1, phone=$2, id_passport_number=$3, nationality=$4, updated_at=$5
		WHERE user_id=$6`
	c.UpdatedAt = time.Now()
	res, err := r.db.Exec(query, c.FullName, c.Phone, c.IDPassportNumber, c.Nationality, c.UpdatedAt, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("guest profile not found")
	}
	return nil
}

func (r *ClientRepository) GetIndividualByID(id uuid.UUID, orgID uuid.UUID) (*models.IndividualClient, error) {
	query := `
		SELECT id, full_name, email, phone, id_passport_number, nationality, status, notes, created_at, updated_at
		FROM individual_profiles
		WHERE id = $1 AND org_id = $2`

	return r.scanIndividual(r.db.QueryRow(query, id, orgID))
}

func (r *ClientRepository) ListIndividual(orgID uuid.UUID, search, status string, page, pageSize int) ([]models.IndividualClient, int, error) {
	where, args, i := r.buildClientWhere(orgID, search, status)

	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM individual_profiles WHERE %s`, where)
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	query := fmt.Sprintf(`
		SELECT id, full_name, email, phone, id_passport_number, nationality, status, notes, created_at, updated_at
		FROM individual_profiles
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, where, i, i+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var clients []models.IndividualClient
	for rows.Next() {
		c, err := r.scanIndividual(rows)
		if err != nil {
			return nil, 0, err
		}
		clients = append(clients, *c)
	}
	return clients, total, rows.Err()
}

func (r *ClientRepository) UpdateIndividual(c *models.IndividualClient, orgID uuid.UUID) error {
	query := `
		UPDATE individual_profiles
		SET full_name=$1, email=$2, phone=$3, id_passport_number=$4,
		    nationality=$5, status=$6, notes=$7, updated_at=$8
		WHERE id=$9 AND org_id=$10`

	c.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		c.FullName, c.Email, c.Phone, c.IDPassportNumber,
		c.Nationality, c.Status, c.Notes, c.UpdatedAt, c.ID, orgID,
	)
	return err
}

func (r *ClientRepository) DeleteIndividual(id uuid.UUID, orgID uuid.UUID) error {
	query := `DELETE FROM individual_profiles WHERE id=$1 AND org_id=$2`

	res, err := r.db.Exec(query, id, orgID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("individual client not found")
	}
	return nil
}

// ─── Corporate ────────────────────────────────────────────────────────────────

func (r *ClientRepository) CreateCorporate(c *models.CorporateClient, orgID uuid.UUID) error {
	query := `
		INSERT INTO corporate_profiles
		    (id, company_name, contact_person, email, phone, company_reg_number, industry, status, notes, org_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

	c.ID = uuid.New()
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	_, err := r.db.Exec(query,
		c.ID, c.CompanyName, c.ContactPerson, c.Email, c.Phone,
		c.CompanyRegNumber, c.Industry, c.Status, c.Notes, orgID, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *ClientRepository) GetCorporateByID(id uuid.UUID, orgID uuid.UUID) (*models.CorporateClient, error) {
	query := `
		SELECT id, company_name, contact_person, email, phone, company_reg_number, industry, status, notes, created_at, updated_at
		FROM corporate_profiles
		WHERE id = $1 AND org_id = $2`

	return r.scanCorporate(r.db.QueryRow(query, id, orgID))
}

func (r *ClientRepository) ListCorporate(orgID uuid.UUID, search, status string, page, pageSize int) ([]models.CorporateClient, int, error) {
	where, args, i := r.buildClientWhere(orgID, search, status)

	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM corporate_profiles WHERE %s`, where)
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	query := fmt.Sprintf(`
		SELECT id, company_name, contact_person, email, phone, company_reg_number, industry, status, notes, created_at, updated_at
		FROM corporate_profiles
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, where, i, i+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var clients []models.CorporateClient
	for rows.Next() {
		c, err := r.scanCorporate(rows)
		if err != nil {
			return nil, 0, err
		}
		clients = append(clients, *c)
	}
	return clients, total, rows.Err()
}

func (r *ClientRepository) UpdateCorporate(c *models.CorporateClient, orgID uuid.UUID) error {
	query := `
		UPDATE corporate_profiles
		SET company_name=$1, contact_person=$2, email=$3, phone=$4,
		    company_reg_number=$5, industry=$6, status=$7, notes=$8, updated_at=$9
		WHERE id=$10 AND org_id=$11`

	c.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		c.CompanyName, c.ContactPerson, c.Email, c.Phone,
		c.CompanyRegNumber, c.Industry, c.Status, c.Notes, c.UpdatedAt, c.ID, orgID,
	)
	return err
}

func (r *ClientRepository) DeleteCorporate(id uuid.UUID, orgID uuid.UUID) error {
	query := `DELETE FROM corporate_profiles WHERE id=$1 AND org_id=$2`

	res, err := r.db.Exec(query, id, orgID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("corporate client not found")
	}
	return nil
}

// LookupIndividualByIDNumber returns the individual client matching an exact
// NRC or passport number within the org — used by the booking dialog pre-flight.
func (r *ClientRepository) LookupIndividualByIDNumber(orgID uuid.UUID, idNumber string) (*models.IndividualClient, error) {
	return r.scanIndividual(r.db.QueryRow(`
		SELECT id, full_name, email, phone, id_passport_number, nationality, status, notes, created_at, updated_at
		FROM individual_profiles
		WHERE org_id = $1 AND id_passport_number = $2`, orgID, idNumber))
}

// SearchCorporate returns corporate clients matching a search term against
// company name or email — used by the booking dialog pre-flight.
func (r *ClientRepository) SearchCorporate(orgID uuid.UUID, search string, limit int) ([]models.CorporateClient, error) {
	fmt.Printf("DEBUG SearchCorporate: org_id=%s search=%q\n", orgID, search)
	rows, err := r.db.Query(`
		SELECT id, company_name, contact_person, email, phone, company_reg_number, industry, status, notes, created_at, updated_at
		FROM corporate_profiles
		WHERE org_id = $1
		  AND (company_name ILIKE $2 OR email ILIKE $2)
		ORDER BY company_name ASC
		LIMIT $3`, orgID, "%"+search+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := []models.CorporateClient{}
	for rows.Next() {
		c, err := r.scanCorporate(rows)
		if err != nil {
			return nil, err
		}
		clients = append(clients, *c)
	}
	fmt.Printf("DEBUG SearchCorporate: found %d result(s)\n", len(clients))
	return clients, rows.Err()
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// buildClientWhere builds a WHERE clause with optional search and status filters.
// Returns the clause string, the args slice, and the next available arg index.
func (r *ClientRepository) buildClientWhere(orgID uuid.UUID, search, status string) (string, []interface{}, int) {
	args := []interface{}{orgID}
	conditions := []string{"org_id = $1"}
	i := 2

	if search != "" {
		conditions = append(conditions, fmt.Sprintf("(email ILIKE $%d OR full_name ILIKE $%d OR company_name ILIKE $%d OR id_passport_number ILIKE $%d)", i, i, i, i))
		args = append(args, "%"+search+"%")
		i++
	}
	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", i))
		args = append(args, status)
		i++
	}

	return strings.Join(conditions, " AND "), args, i
}

func (r *ClientRepository) scanIndividual(row rowScanner) (*models.IndividualClient, error) {
	var c models.IndividualClient
	var idPassport, nationality, notes sql.NullString

	err := row.Scan(
		&c.ID, &c.FullName, &c.Email, &c.Phone, &idPassport,
		&nationality, &c.Status, &notes, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if idPassport.Valid {
		c.IDPassportNumber = idPassport.String
	}
	if nationality.Valid {
		c.Nationality = nationality.String
	}
	if notes.Valid {
		c.Notes = notes.String
	}
	return &c, nil
}

func (r *ClientRepository) scanCorporate(row rowScanner) (*models.CorporateClient, error) {
	var c models.CorporateClient
	var industry, notes sql.NullString

	err := row.Scan(
		&c.ID, &c.CompanyName, &c.ContactPerson, &c.Email, &c.Phone,
		&c.CompanyRegNumber, &industry, &c.Status, &notes, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if industry.Valid {
		c.Industry = industry.String
	}
	if notes.Valid {
		c.Notes = notes.String
	}
	return &c, nil
}
