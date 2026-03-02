package dto

import "censys/pkg/models"

// CreateAssetRequest represents the request body for the CreateAsset endpoint
type CreateAssetRequest struct {
	IPAddress   string   `json:"ip_address" binding:"required"`
	Hostname    string   `json:"hostname" binding:"required"`
	PortNumbers []int    `json:"port_numbers" binding:"required"`
	Tags        []string `json:"tags"`
}

// CreateAssetResponse represents the response body for the CreateAsset endpoint
type CreateAssetResponse struct {
	ID          string   `json:"id"`
	IPAddress   string   `json:"ip_address"`
	Hostname    string   `json:"hostname"`
	RiskLevel   string   `json:"risk_level"`
	PortNumbers []int    `json:"port_numbers"`
	Tags        []string `json:"tags"`
}

type ListAssetsRequest struct {
	Search    string `form:"search"`
	RiskLevel string `form:"risk_level"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}

type ListAssetsResponse struct {
	Assets     []models.Asset `json:"assets"`
	Pagination Pagination     `json:"pagination"`
}
