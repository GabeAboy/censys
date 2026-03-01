package dto

type CreateAssetTagRequest struct {
	TagName string `json:"tag_name" binding:"required"`
}
