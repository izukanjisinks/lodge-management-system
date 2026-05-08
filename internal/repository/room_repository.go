package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{db: database.DB}
}

func (r *RoomRepository) Create(room *models.Room, orgID uuid.UUID) error {
	room.ID = uuid.New()
	now := time.Now()
	room.CreatedAt = now
	room.UpdatedAt = now
	_, err := r.db.Exec(`
		INSERT INTO rooms (id, name, type, capacity, price_per_night, amenities, images, is_available, description, org_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		room.ID, room.Name, room.Type, room.Capacity, room.PricePerNight,
		pq.Array(room.Amenities), pq.Array(room.Images), room.IsAvailable, room.Description,
		orgID, room.CreatedAt, room.UpdatedAt,
	)
	return err
}

// GetByIDUnscoped fetches a room by ID with no org filter — use only in guest flows
// where org is derived from the room itself after lookup.
func (r *RoomRepository) GetByIDUnscoped(id uuid.UUID) (*models.Room, error) {
	row := r.db.QueryRow(`
		SELECT id, org_id, name, type, capacity, price_per_night, amenities, images, is_available, description, created_at, updated_at
		FROM rooms WHERE id = $1`, id)
	return scanRoom(row)
}

func (r *RoomRepository) GetByID(id uuid.UUID, orgID uuid.UUID) (*models.Room, error) {
	row := r.db.QueryRow(`
		SELECT id, org_id, name, type, capacity, price_per_night, amenities, images, is_available, description, created_at, updated_at
		FROM rooms WHERE id = $1 AND org_id = $2`, id, orgID)
	return scanRoom(row)
}

// GuestList lists rooms with optional filters — used by public guest endpoints.
func (r *RoomRepository) GuestList(orgID *uuid.UUID, roomType, orgName string, isAvailable *bool, page, pageSize int) ([]models.Room, int, error) {
	args := []interface{}{}
	where := []string{"1=1"}
	i := 1

	if orgID != nil {
		where = append(where, fmt.Sprintf("r.org_id = $%d", i))
		args = append(args, *orgID)
		i++
	}
	if roomType != "" {
		where = append(where, fmt.Sprintf("r.type = $%d", i))
		args = append(args, roomType)
		i++
	}
	if isAvailable != nil {
		where = append(where, fmt.Sprintf("r.is_available = $%d", i))
		args = append(args, *isAvailable)
		i++
	}
	if orgName != "" {
		where = append(where, fmt.Sprintf("o.name ILIKE $%d", i))
		args = append(args, "%"+orgName+"%")
		i++
	}

	whereStr := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM rooms r LEFT JOIN organizations o ON o.id = r.org_id WHERE %s`, whereStr), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT r.id, r.org_id, r.name, r.type, r.capacity, r.price_per_night, r.amenities, r.images, r.is_available, r.description, r.created_at, r.updated_at,
		       o.name, o.email, o.address, o.phone, o.logo_url
		FROM rooms r
		LEFT JOIN organizations o ON o.id = r.org_id
		WHERE %s
		ORDER BY r.name ASC
		LIMIT $%d OFFSET $%d`, whereStr, i, i+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		var orgID sql.NullString
		var description sql.NullString
		var orgName, orgEmail, orgAddress, orgPhone, orgLogoURL sql.NullString
		if err := rows.Scan(
			&room.ID, &orgID, &room.Name, &room.Type, &room.Capacity, &room.PricePerNight,
			pq.Array(&room.Amenities), pq.Array(&room.Images), &room.IsAvailable, &description,
			&room.CreatedAt, &room.UpdatedAt,
			&orgName, &orgEmail, &orgAddress, &orgPhone, &orgLogoURL,
		); err != nil {
			return nil, 0, err
		}
		if orgID.Valid {
			parsed, _ := uuid.Parse(orgID.String)
			room.OrgID = &parsed
		}
		if description.Valid {
			room.Description = description.String
		}
		if room.Amenities == nil {
			room.Amenities = []string{}
		}
		if room.Images == nil {
			room.Images = []string{}
		}
		if orgName.Valid {
			room.Organization = &models.RoomOrganization{
				Name:    orgName.String,
				Email:   orgEmail.String,
				Address: orgAddress.String,
				Phone:   orgPhone.String,
				LogoURL: orgLogoURL.String,
			}
		}
		rooms = append(rooms, room)
	}
	return rooms, total, rows.Err()
}

func (r *RoomRepository) List(orgID uuid.UUID, roomType string, isAvailable *bool, page, pageSize int) ([]models.Room, int, error) {
	args := []interface{}{orgID}
	where := []string{"org_id = $1"}
	i := 2

	if roomType != "" {
		where = append(where, fmt.Sprintf("type = $%d", i))
		args = append(args, roomType)
		i++
	}
	if isAvailable != nil {
		where = append(where, fmt.Sprintf("is_available = $%d", i))
		args = append(args, *isAvailable)
		i++
	}

	whereStr := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM rooms WHERE %s`, whereStr), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT id, org_id, name, type, capacity, price_per_night, amenities, images, is_available, description, created_at, updated_at
		FROM rooms WHERE %s
		ORDER BY name ASC
		LIMIT $%d OFFSET $%d`, whereStr, i, i+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		room, err := scanRoom(rows)
		if err != nil {
			return nil, 0, err
		}
		rooms = append(rooms, *room)
	}
	return rooms, total, rows.Err()
}

func (r *RoomRepository) ListAvailable(orgID uuid.UUID, checkIn, checkOut time.Time, roomType string) ([]models.Room, error) {
	args := []interface{}{checkOut, checkIn, orgID}
	extra := ""
	if roomType != "" {
		args = append(args, roomType)
		extra = fmt.Sprintf(" AND type = $%d", len(args))
	}

	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT id, org_id, name, type, capacity, price_per_night, amenities, images, is_available, description, created_at, updated_at
		FROM rooms
		WHERE is_available = TRUE AND org_id = $3%s
		  AND id NOT IN (
		    SELECT room_id FROM bookings
		    WHERE status IN ('confirmed', 'checked_in')
		      AND check_in  < $1
		      AND check_out > $2
		  )
		ORDER BY name ASC`, extra), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		room, err := scanRoom(rows)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, *room)
	}
	return rooms, rows.Err()
}

func (r *RoomRepository) Update(room *models.Room, orgID uuid.UUID) error {
	room.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE rooms SET name=$1, type=$2, capacity=$3, price_per_night=$4,
		       amenities=$5, is_available=$6, description=$7, updated_at=$8
		WHERE id=$9 AND org_id=$10`,
		room.Name, room.Type, room.Capacity, room.PricePerNight,
		pq.Array(room.Amenities), room.IsAvailable, room.Description,
		room.UpdatedAt, room.ID, orgID,
	)
	return err
}

func (r *RoomRepository) UpdateImages(id uuid.UUID, orgID uuid.UUID, images []string) error {
	_, err := r.db.Exec(`UPDATE rooms SET images=$1, updated_at=$2 WHERE id=$3 AND org_id=$4`,
		pq.Array(images), time.Now(), id, orgID)
	return err
}

func (r *RoomRepository) SetAvailability(id uuid.UUID, orgID uuid.UUID, available bool) error {
	_, err := r.db.Exec(`UPDATE rooms SET is_available=$1, updated_at=$2 WHERE id=$3 AND org_id=$4`,
		available, time.Now(), id, orgID)
	return err
}

func (r *RoomRepository) Delete(id uuid.UUID, orgID uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM rooms WHERE id=$1 AND org_id=$2`, id, orgID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("room not found")
	}
	return nil
}

type roomScanner interface {
	Scan(dest ...interface{}) error
}

func scanRoom(row roomScanner) (*models.Room, error) {
	var room models.Room
	var orgID sql.NullString
	var description sql.NullString
	err := row.Scan(
		&room.ID, &orgID, &room.Name, &room.Type, &room.Capacity, &room.PricePerNight,
		pq.Array(&room.Amenities), pq.Array(&room.Images), &room.IsAvailable, &description,
		&room.CreatedAt, &room.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if orgID.Valid {
		parsed, _ := uuid.Parse(orgID.String)
		room.OrgID = &parsed
	}
	if description.Valid {
		room.Description = description.String
	}
	if room.Amenities == nil {
		room.Amenities = []string{}
	}
	if room.Images == nil {
		room.Images = []string{}
	}
	return &room, nil
}
