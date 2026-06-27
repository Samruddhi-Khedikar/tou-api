package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samkhedikar/tou-api/internal/model"
)

type ChargerRepository interface {
	Insert(ctx context.Context, charger model.Charger) error
	FindByID(ctx context.Context, id uuid.UUID) (model.Charger, error)
}

type chargerRepo struct {
	db *sql.DB
}

func NewChargerRepo(db *sql.DB) ChargerRepository {
	return &chargerRepo{db: db}
}

func (r *chargerRepo) Insert(ctx context.Context, charger model.Charger) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO chargers (id, name, location, timezone) VALUES (?, ?, ?, ?)`,
		charger.ID.String(), charger.Name, charger.Location, charger.Timezone,
	)
	if err != nil {
		return fmt.Errorf("insert charger: %w", err)
	}
	return nil
}

func (r *chargerRepo) FindByID(ctx context.Context, id uuid.UUID) (model.Charger, error) {
	var c model.Charger
	var createdAt, updatedAt string
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, location, timezone, created_at, updated_at FROM chargers WHERE id = ?`,
		id.String(),
	).Scan(&c.ID, &c.Name, &c.Location, &c.Timezone, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return model.Charger{}, ErrNotFound
	}
	if err != nil {
		return model.Charger{}, fmt.Errorf("find charger: %w", err)
	}
	c.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	c.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return c, nil
}
