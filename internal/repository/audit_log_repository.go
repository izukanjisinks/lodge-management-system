package repository

import (
	"database/sql"
	"fmt"
	"time"

	"lodge-system/internal/database"
	"lodge-system/internal/models"

	"github.com/google/uuid"
)

type AuditLogRepository struct {
	db *sql.DB
}

func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{db: database.DB}
}

func (r *AuditLogRepository) Insert(log *models.AuditLog) error {
	log.ID = uuid.New()
	log.CreatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO audit_logs (id, org_id, actor_type, actor_id, actor_name, action, entity_type, entity_id, payload, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		log.ID, log.OrgID, log.ActorType, log.ActorID, log.ActorName,
		log.Action, log.EntityType, log.EntityID, log.Payload, log.CreatedAt,
	)
	return err
}

func (r *AuditLogRepository) List(orgID uuid.UUID, entityType, entityID, action string, page, pageSize int) ([]models.AuditLog, int, error) {
	args := []interface{}{orgID}
	where := []string{"org_id = $1"}
	i := 2

	if entityType != "" {
		where = append(where, fmt.Sprintf("entity_type = $%d", i))
		args = append(args, entityType)
		i++
	}
	if entityID != "" {
		id, err := uuid.Parse(entityID)
		if err == nil {
			where = append(where, fmt.Sprintf("entity_id = $%d", i))
			args = append(args, id)
			i++
		}
	}
	if action != "" {
		where = append(where, fmt.Sprintf("action = $%d", i))
		args = append(args, action)
		i++
	}

	whereStr := fmt.Sprintf("WHERE %s", where[0])
	for _, w := range where[1:] {
		whereStr += " AND " + w
	}

	var total int
	if err := r.db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM audit_logs %s`, whereStr), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pageSize, (page-1)*pageSize)
	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT id, org_id, actor_type, actor_id, actor_name, action, entity_type, entity_id, payload, created_at
		FROM audit_logs %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereStr, i, i+1), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		var actorID uuid.NullUUID
		var actorName sql.NullString
		if err := rows.Scan(&l.ID, &l.OrgID, &l.ActorType, &actorID, &actorName,
			&l.Action, &l.EntityType, &l.EntityID, &l.Payload, &l.CreatedAt); err != nil {
			return nil, 0, err
		}
		if actorID.Valid {
			l.ActorID = &actorID.UUID
		}
		if actorName.Valid {
			l.ActorName = actorName.String
		}
		logs = append(logs, l)
	}
	if logs == nil {
		logs = []models.AuditLog{}
	}
	return logs, total, rows.Err()
}
