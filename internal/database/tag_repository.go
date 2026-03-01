package database

import (
	"database/sql"
	"fmt"

	"censys/pkg/models"

	"github.com/lib/pq"
)

type TagRepository struct {
	db *DB
}

func NewTagRepository(db *DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(assetID string, tag models.Tag) error {
	query := `
		INSERT INTO tags (asset_id, tag_name, created_at)
		VALUES ($1, $2, NOW())
	`

	_, err := r.db.Exec(
		query,
		assetID,
		tag.TagName,
	)

	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}

func (r *TagRepository) BulkCreateWithTx(tx *sql.Tx, assetID *models.AssetId, tags []*string) error {
	query := `
		INSERT INTO tags (asset_id, tag_name, created_at)
		VALUES ($1, $2, NOW())
	`

	for _, tagName := range tags {
		_, err := tx.Exec(query, assetID, tagName)
		if err != nil {
			pqErr, _ := err.(*pq.Error)
			if pqErr.Constraint == models.UniqueAssetTag {
				continue
			}
			return fmt.Errorf("failed to create tag '%s': %w", tagName, err)
		}
	}

	return nil
}
