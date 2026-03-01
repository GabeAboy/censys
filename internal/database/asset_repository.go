package database

import (
	"database/sql"
	"fmt"

	"censys/pkg/models"

	"github.com/lib/pq"
)

type AssetRepository struct {
	db *DB
}

func NewAssetRepository(db *DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) Create(tx *sql.Tx, asset *models.Asset) (*models.AssetId, error) {
	query := `
		INSERT INTO assets (ip_address, hostname, risk_level, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id`

	var assetId models.AssetId
	err := tx.QueryRow(query, asset.IPAddress, asset.Hostname, asset.RiskLevel).Scan(&assetId)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}

	return &assetId, nil
}

// TODO This query should also include ports that aren't in the filter set
func (r *AssetRepository) GetAll(searchText *string, riskLevels []*models.RiskLevel, limit int, offset int) ([]models.Asset, error) {
	query := `
		SELECT
			a.id,
			a.ip_address,
			a.hostname,
			a.risk_level,
			a.created_at,
			COALESCE(array_agg(DISTINCT p.port_number) FILTER (WHERE p.port_number IS NOT NULL), '{}') as ports,
			COALESCE(array_agg(DISTINCT t.tag_name) FILTER (WHERE t.tag_name IS NOT NULL), '{}') as tags
		FROM assets a
		LEFT JOIN ports p ON a.id = p.asset_id
		LEFT JOIN tags t ON a.id = t.asset_id
		WHERE (
			a.hostname ILIKE CONCAT('%', COALESCE($1, ''), '%') OR
			a.ip_address ILIKE CONCAT('%', COALESCE($1, ''), '%') OR
			CAST(p.port_number AS TEXT) LIKE CONCAT('%', COALESCE($1, ''), '%') OR
			t.tag_name ILIKE CONCAT('%', COALESCE($1, ''), '%')
		)
		AND a.risk_level = ANY($2)
		GROUP BY a.id, a.ip_address, a.hostname, a.risk_level, a.created_at
		ORDER BY a.risk_level, a.created_at DESC, a.hostname, a.ip_address DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, searchText, riskLevels, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query assets: %w", err)
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var asset models.Asset

		err := rows.Scan(
			&asset.ID,
			&asset.IPAddress,
			&asset.Hostname,
			&asset.RiskLevel,
			&asset.CreatedAt,
			pq.Array(&asset.Ports),
			pq.Array(&asset.Tags),
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

func (r *AssetRepository) GetAllCount(searchText *string, riskLevels []*models.RiskLevel) (int, error) {
	query := `
		SELECT COUNT(DISTINCT a.id)
		FROM assets a
		LEFT JOIN ports p ON a.id = p.asset_id
		LEFT JOIN tags t ON a.id = t.asset_id
		WHERE (
			a.hostname ILIKE CONCAT('%', COALESCE($1, ''), '%') OR
			a.ip_address ILIKE CONCAT('%', COALESCE($1, ''), '%') OR
			CAST(p.port_number AS TEXT) LIKE CONCAT('%', COALESCE($1, ''), '%') OR
			t.tag_name ILIKE CONCAT('%', COALESCE($1, ''), '%')
		)
		AND a.risk_level = ANY($2)
	`

	var count int
	err := r.db.QueryRow(query, searchText, riskLevels).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count assets: %w", err)
	}

	return count, nil
}

func (r *AssetRepository) Delete(id string) error {
	query := `DELETE FROM assets WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	return nil
}
