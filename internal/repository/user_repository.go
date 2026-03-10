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

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.DB}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (user_id, email, password, role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	user.UserID = uuid.New()
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	_, err := r.db.Exec(query,
		user.UserID, user.Email, user.Password, user.RoleID, user.IsActive, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	query := `
		SELECT u.user_id, u.email, u.password, u.role_id, u.is_active, u.created_at, u.updated_at,
		       u.change_password, u.password_changed_at, u.password_expires_at,
		       u.failed_login_attempts, u.is_locked, u.locked_until,
		       r.role_id, r.name, r.description
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		WHERE u.user_id = $1`

	return r.scanUser(r.db.QueryRow(query, id))
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT u.user_id, u.email, u.password, u.role_id, u.is_active, u.created_at, u.updated_at,
		       u.change_password, u.password_changed_at, u.password_expires_at,
		       u.failed_login_attempts, u.is_locked, u.locked_until,
		       r.role_id, r.name, r.description
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		WHERE u.email = $1`

	return r.scanUser(r.db.QueryRow(query, email))
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query := `
		SELECT u.user_id, u.email, u.password, u.role_id, u.is_active, u.created_at, u.updated_at,
		       u.change_password, u.password_changed_at, u.password_expires_at,
		       u.failed_login_attempts, u.is_locked, u.locked_until,
		       r.role_id, r.name, r.description
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		ORDER BY u.created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		u, err := r.scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, rows.Err()
}

// List returns paginated users with optional filtering
func (r *UserRepository) List(search string, roleID *uuid.UUID, isActive *bool, page, pageSize int) ([]models.User, int, error) {
	args := []interface{}{}
	where := []string{}
	i := 1

	// Search filter
	if search != "" {
		where = append(where, fmt.Sprintf("(u.email ILIKE $%d)", i))
		args = append(args, "%"+search+"%")
		i++
	}

	// Role filter
	if roleID != nil {
		where = append(where, fmt.Sprintf("u.role_id=$%d", i))
		args = append(args, *roleID)
		i++
	}

	// Active status filter
	if isActive != nil {
		where = append(where, fmt.Sprintf("u.is_active=$%d", i))
		args = append(args, *isActive)
		i++
	}

	whereStr := "1=1"
	if len(where) > 0 {
		whereStr = strings.Join(where, " AND ")
	}

	// Get total count
	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM users u WHERE %s`, whereStr)
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	args = append(args, pageSize, (page-1)*pageSize)
	query := fmt.Sprintf(`
		SELECT u.user_id, u.email, u.password, u.role_id, u.is_active, u.created_at, u.updated_at,
		       u.change_password, u.password_changed_at, u.password_expires_at,
		       u.failed_login_attempts, u.is_locked, u.locked_until,
		       r.role_id, r.name, r.description
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.role_id
		WHERE %s
		ORDER BY u.created_at DESC
		LIMIT $%d OFFSET $%d`, whereStr, i, i+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		u, err := r.scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, *u)
	}

	return users, total, rows.Err()
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users SET email = $1, role_id = $2, is_active = $3, updated_at = $4,
		       password = $5, change_password = $6, password_changed_at = $7,
		       password_expires_at = $8, failed_login_attempts = $9,
		       is_locked = $10, locked_until = $11
		WHERE user_id = $12`
	_, err := r.db.Exec(query,
		user.Email, user.RoleID, user.IsActive, time.Now(),
		user.Password, user.ChangePassword, user.PasswordChangedAt,
		user.PasswordExpiresAt, user.FailedLoginAttempts,
		user.IsLocked, user.LockedUntil, user.UserID,
	)
	return err
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE user_id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(1) FROM users WHERE email = $1`, email).Scan(&count)
	return count > 0, err
}

// scanUser works with both *sql.Row and *sql.Rows via a common interface
type rowScanner interface {
	Scan(dest ...interface{}) error
}

func (r *UserRepository) scanUser(row rowScanner) (*models.User, error) {
	var u models.User
	var roleID sql.NullString
	var rRoleID, rName, rDesc sql.NullString
	var passwordChangedAt, passwordExpiresAt, lockedUntil sql.NullTime

	err := row.Scan(
		&u.UserID, &u.Email, &u.Password, &roleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
		&u.ChangePassword, &passwordChangedAt, &passwordExpiresAt,
		&u.FailedLoginAttempts, &u.IsLocked, &lockedUntil,
		&rRoleID, &rName, &rDesc,
	)
	if err != nil {
		return nil, err
	}

	if roleID.Valid {
		parsed, _ := uuid.Parse(roleID.String)
		u.RoleID = &parsed
	}

	if passwordChangedAt.Valid {
		u.PasswordChangedAt = &passwordChangedAt.Time
	}

	if passwordExpiresAt.Valid {
		u.PasswordExpiresAt = &passwordExpiresAt.Time
	}

	if lockedUntil.Valid {
		u.LockedUntil = &lockedUntil.Time
	}

	if rRoleID.Valid {
		roleUUID, _ := uuid.Parse(rRoleID.String)
		u.Role = &models.Role{
			RoleID:      roleUUID,
			Name:        rName.String,
			Description: rDesc.String,
		}
	}

	return &u, nil
}

// GetUserWithFewestTasksByRole finds a user with the specified role who has the fewest pending tasks
// This enables load balancing when assigning workflow tasks
func (r *UserRepository) GetUserWithFewestTasksByRole(roleName string) (*models.User, error) {
	query := `
		SELECT u.user_id, u.email, u.password, u.role_id, u.is_active, u.created_at, u.updated_at,
		       r.role_id, r.name, r.description,
		       COALESCE(COUNT(at.id), 0) as task_count
		FROM users u
		INNER JOIN roles r ON u.role_id = r.role_id
		LEFT JOIN assigned_tasks at ON at.assigned_to = u.user_id AND at.status = 'pending'
		WHERE r.name = $1 AND u.is_active = true
		GROUP BY u.user_id, u.email, u.password, u.role_id, u.is_active, u.created_at, u.updated_at,
		         r.role_id, r.name, r.description
		ORDER BY task_count ASC, u.created_at ASC
		LIMIT 1
	`

	var u models.User
	var roleID sql.NullString
	var rRoleID, rName, rDesc sql.NullString
	var taskCount int

	err := r.db.QueryRow(query, roleName).Scan(
		&u.UserID, &u.Email, &u.Password, &roleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
		&rRoleID, &rName, &rDesc, &taskCount,
	)
	if err != nil {
		return nil, err
	}

	if roleID.Valid {
		parsed, _ := uuid.Parse(roleID.String)
		u.RoleID = &parsed
	}

	if rRoleID.Valid {
		roleUUID, _ := uuid.Parse(rRoleID.String)
		u.Role = &models.Role{
			RoleID:      roleUUID,
			Name:        rName.String,
			Description: rDesc.String,
		}
	}

	return &u, nil
}
