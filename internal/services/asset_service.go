package services

import (
	"censys/internal/database"
	"censys/internal/handlers/dto"
	"censys/pkg/models"
	"context"
	"fmt"
)

type AssetService struct {
	db        *database.DB
	assetRepo *database.AssetRepository
	portRepo  *database.PortRepository
	tagRepo   *database.TagRepository
}

func NewAssetService(
	db *database.DB,
	assetRepo *database.AssetRepository,
	portRepo *database.PortRepository,
	tagRepo *database.TagRepository,
) *AssetService {
	return &AssetService{
		db:        db,
		assetRepo: assetRepo,
		portRepo:  portRepo,
		tagRepo:   tagRepo,
	}
}

func (s *AssetService) CreateAsset(ctx context.Context, assetReq dto.CreateAssetRequest) error {
	asset := &models.Asset{
		IPAddress: assetReq.IPAddress,
		Hostname:  assetReq.Hostname,
	}

	ports, tags := getUniquePortsAndTags(assetReq.PortNumbers, assetReq.Tags)
	asset.RiskLevel = calculateRiskLevel(ports)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	assetId, err := s.assetRepo.Create(ctx, tx, asset)
	if err != nil || assetId == nil {
		return fmt.Errorf("failed to create asset: %w", err)
	}

	if len(ports) > 0 {
		if err := s.portRepo.Create(ctx, tx, assetId, ports); err != nil {
			return fmt.Errorf("failed to create ports: %w", err)
		}
	}

	// Only create tags if there are any
	if assetReq.Tags != nil && len(tags) > 0 {
		if err := s.tagRepo.BulkCreateWithTx(ctx, tx, assetId, tags); err != nil {
			return fmt.Errorf("failed to create tags: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *AssetService) GetAssetList(ctx context.Context, searchText string, riskLevels []models.RiskLevel, offset int, limit int) ([]models.Asset, int, error) {
	asset, err := s.assetRepo.GetAll(ctx, searchText, riskLevels, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get asset: %w", err)
	}

	totalAssets, err := s.assetRepo.GetAllCount(ctx, searchText, riskLevels)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get asset counts: %w", err)
	}
	return asset, totalAssets, nil
}

func (s *AssetService) GetAssetCount(ctx context.Context, searchText string, riskLevels []models.RiskLevel) (int, error) {
	totalAssets, err := s.assetRepo.GetAllCount(ctx, searchText, riskLevels)
	if err != nil {
		return 0, fmt.Errorf("failed to get asset count: %w", err)
	}
	return totalAssets, nil
}

func (s *AssetService) DeleteAsset(ctx context.Context, assetID string) error {
	if err := s.assetRepo.Delete(ctx, assetID); err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}

	return nil
}

func (s *AssetService) CreateAssetTag(ctx context.Context, assetID string, createTagReq dto.CreateAssetTagRequest) error {
	tag := models.Tag{
		TagName: &createTagReq.TagName,
	}

	if err := s.tagRepo.Create(ctx, assetID, tag); err != nil {
		return fmt.Errorf("failed to add tag: %w", err)
	}

	return nil
}

func calculateRiskLevel(ports []models.Port) models.RiskLevel {
	hasHighRisk := false
	hasMediumRisk := false

	for _, port := range ports {
		if port.PortNumber == 22 || port.PortNumber == 3389 || port.PortNumber == 21 {
			hasHighRisk = true
		}

		if port.PortNumber == 443 && port.HasExpiredCert {
			hasMediumRisk = true
		}
	}

	if hasHighRisk {
		return models.RiskLevelHigh
	}
	if hasMediumRisk {
		return models.RiskLevelMedium
	}
	return models.RiskLevelLow
}

func getUniquePortsAndTags(portNumbers []int, tagNames []string) ([]models.Port, []string) {
	portMap := make(map[int]struct{})
	for _, portNum := range portNumbers {
		portMap[portNum] = struct{}{}
	}

	ports := make([]models.Port, 0, len(portMap))
	for portNum := range portMap {
		port := models.Port{
			PortNumber:     portNum,
			HasExpiredCert: false,
		}
		ports = append(ports, port)
	}

	tagsMap := make(map[string]struct{})
	for _, tagName := range tagNames {
		if tagName != "" {
			tagsMap[tagName] = struct{}{}
		}
	}

	tags := make([]string, 0, len(tagsMap))
	for tagStr := range tagsMap {
		tags = append(tags, tagStr)
	}

	return ports, tags
}
