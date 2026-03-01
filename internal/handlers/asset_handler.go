package handlers

import (
	"censys/internal/handlers/dto"
	"censys/pkg/models"
	"net/http"
	"strings"

	"censys/internal/database"
	"censys/internal/services"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	assetRepo    *database.AssetRepository
	assetService *services.AssetService
}

func NewAssetHandler(assetRepo *database.AssetRepository, assetService *services.AssetService) *AssetHandler {
	return &AssetHandler{
		assetRepo:    assetRepo,
		assetService: assetService,
	}
}

// GetAssetList godoc
// @Summary List assets
// @Description Get a paginated list of assets with optional search and risk level filters
// @Tags assets
// @Accept json
// @Produce json
// @Param search query string false "Search term for hostname, IP, ports, or tags"
// @Param risk_level query string false "Comma-separated risk levels (e.g., High,Medium,Low)"
// @Param page query int false "Page number (offset)" default(0)
// @Param page_size query int false "Items per page (limit)" default(10)
// @Success 200 {object} dto.ListAssetsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /assets [get]
func (h *AssetHandler) GetAssetList(c *gin.Context) {
	var req dto.ListAssetsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse comma-separated risk levels
	riskLevels := parseRiskLevels(req.RiskLevel)
	p := req.Page - 1
	offset := p * req.PageSize
	assets, totalAssets, err := h.assetService.GetAssetList(req.Search, riskLevels, offset, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve assets"})
		return
	}

	pageSize := len(*assets)
	response := dto.ListAssetsResponse{
		Assets: assets,
		Pagination: dto.Pagination{
			Total:    totalAssets,
			Page:     req.Page,
			PageSize: pageSize,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetAssetCount godoc
// @Summary Get asset count
// @Description Get the total count of assets matching the search and filter criteria
// @Tags assets
// @Accept json
// @Produce json
// @Param search query string false "Search term for hostname, IP, ports, or tags"
// @Param risk_level query string false "Comma-separated risk levels (e.g., High,Medium,Low)"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /assets/count [get]
func (h *AssetHandler) GetAssetCount(c *gin.Context) {
	var req dto.ListAssetsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse comma-separated risk levels
	riskLevels := parseRiskLevels(req.RiskLevel)

	totalAssets, err := h.assetService.GetAssetCount(req.Search, riskLevels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve asset count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": totalAssets})
}

// CreateAsset godoc
// @Summary Create a new asset
// @Description Create a new asset with ports and tags, risk level is automatically calculated
// @Tags assets
// @Accept json
// @Produce json
// @Param asset body dto.CreateAssetRequest true "Asset creation data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /assets [post]
func (h *AssetHandler) CreateAsset(c *gin.Context) {
	var req dto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.assetService.CreateAsset(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create asset"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Asset created successfully"})
}

// DeleteAsset godoc
// @Summary Delete an asset
// @Description Delete an asset by ID (cascades to ports and tags)
// @Tags assets
// @Accept json
// @Produce json
// @Param id path string true "Asset ID (UUID)"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /assets/{id} [delete]
func (h *AssetHandler) DeleteAsset(c *gin.Context) {
	assetID := c.Param("id")

	if err := h.assetService.DeleteAsset(assetID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete asset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset deleted successfully"})
}

// CreateAssetTag godoc
// @Summary Add a tag to an asset
// @Description Add a new tag to an existing asset
// @Tags assets
// @Accept json
// @Produce json
// @Param id path string true "Asset ID (UUID)"
// @Param tag body dto.CreateAssetTagRequest true "Tag data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /assets/{id}/tags [post]
func (h *AssetHandler) CreateAssetTag(c *gin.Context) {
	assetID := c.Param("id")

	var req dto.CreateAssetTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.assetService.CreateAssetTag(assetID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tag"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tag added successfully"})
}

// parseRiskLevels parses a comma-separated string of risk levels into a slice
func parseRiskLevels(riskLevelStr string) []*models.RiskLevel {
	if riskLevelStr == "" {
		return nil
	}

	parts := strings.Split(riskLevelStr, ",")
	riskLevels := make([]*models.RiskLevel, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			riskLevel := models.RiskLevel(trimmed)
			riskLevels = append(riskLevels, &riskLevel)
		}
	}

	return riskLevels
}
