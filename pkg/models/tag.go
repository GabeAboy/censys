package models

import (
	"time"
)

type Tag struct {
	ID        *string    `json:"id" db:"id"`
	AssetID   *string    `json:"asset_id" db:"asset_id"`
	TagName   *string    `json:"tag_name" db:"tag_name"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
}
