package models

type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "Low"
	RiskLevelMedium RiskLevel = "Medium"
	RiskLevelHigh   RiskLevel = "High"
)

type AssetId string

const (
	UniqueAssetPort string = "unique_asset_port"
	UniqueAssetTag  string = "unique_asset_tag"
)
