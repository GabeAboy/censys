package models

import (
	"time"
)

type Port struct {
	ID             string    `json:"id" db:"id"`
	AssetID        string    `json:"asset_id" db:"asset_id"`
	PortNumber     int       `json:"port_number" db:"port_number"`
	HasExpiredCert bool      `json:"has_expired_cert" db:"has_expired_cert"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
