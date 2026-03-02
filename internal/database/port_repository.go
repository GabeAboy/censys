package database

import (
	"censys/pkg/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type PortRepository struct {
	db *DB
}

func NewPortRepository(db *DB) *PortRepository {
	return &PortRepository{db: db}
}

func (r *PortRepository) Create(ctx context.Context, tx *sql.Tx, assetID *models.AssetId, ports []models.Port) error {
	query := `
			INSERT INTO ports (asset_id, port_number, created_at)
			VALUES ($1, $2, NOW())
	`

	for _, port := range ports {
		_, err := tx.ExecContext(ctx, query, assetID, port.PortNumber)
		if err != nil {
			pqErr, _ := err.(*pq.Error)
			if pqErr.Constraint == models.UniqueAssetPort {
				continue
			}
			return fmt.Errorf("failed to create port: %w", err)
		}
	}

	return nil
}
