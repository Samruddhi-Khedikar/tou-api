package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samkhedikar/tou-api/internal/model"
)

type PricingRepository interface {
	// ReplaceAll deletes all existing periods for a charger and inserts the new ones atomically.
	ReplaceAll(ctx context.Context, chargerID uuid.UUID, periods []model.PricingPeriod) error
	FindByChargerID(ctx context.Context, chargerID uuid.UUID) ([]model.PricingPeriod, error)
}

type pricingRepo struct {
	db *sql.DB
}

func NewPricingRepo(db *sql.DB) PricingRepository {
	return &pricingRepo{db: db}
}

// ReplaceAll replaces all pricing periods for a charger in a single transaction.
// Deletes old periods first, then inserts new ones.
// If any insert fails, the transaction rolls back and old periods are preserved.
func (r *pricingRepo) ReplaceAll(ctx context.Context, chargerID uuid.UUID, periods []model.PricingPeriod) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
		return err
	}

	// Delete all existing periods for this charger.
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM pricing_periods WHERE charger_id = ?`,
		chargerID.String(),
	); err != nil {
		return fmt.Errorf("delete existing periods: %w", err)
	}

	// Insert new periods.
	for _, p := range periods {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO pricing_periods (id, charger_id, start_hour, start_minute, end_hour, end_minute, price_per_kwh)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			p.ID.String(), chargerID.String(),
			p.StartHour, p.StartMinute,
			p.EndHour, p.EndMinute,
			p.PricePerKWh,
		); err != nil {
			return fmt.Errorf("insert period: %w", err)
		}
	}

	return tx.Commit()
}

func (r *pricingRepo) FindByChargerID(ctx context.Context, chargerID uuid.UUID) ([]model.PricingPeriod, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, charger_id, start_hour, start_minute, end_hour, end_minute, price_per_kwh, created_at
		 FROM pricing_periods
		 WHERE charger_id = ?
		 ORDER BY start_hour, start_minute`,
		chargerID.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("find periods: %w", err)
	}
	defer rows.Close()

	var periods []model.PricingPeriod
	for rows.Next() {
		var p model.PricingPeriod
		var createdAt string
		if err := rows.Scan(&p.ID, &p.ChargerID, &p.StartHour, &p.StartMinute,
			&p.EndHour, &p.EndMinute, &p.PricePerKWh, &createdAt); err != nil {
			return nil, fmt.Errorf("scan period: %w", err)
		}
		p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		periods = append(periods, p)
	}
	return periods, nil
}
