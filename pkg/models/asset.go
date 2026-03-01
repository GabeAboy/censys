package models

import (
	"time"
)

type Asset struct {
	ID        string    `json:"id" db:"id"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	Hostname  string    `json:"hostname" db:"hostname"`
	RiskLevel RiskLevel `json:"risk_level" db:"risk_level"`
	Ports     []int64   `json:"ports" db:"ports"`
	Tags      []string  `json:"tags" db:"tags"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
