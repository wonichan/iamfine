package repository

import (
	"context"

	"hupu/shared/models"
)

type LikeRepository interface {
	Like(ctx context.Context, userID, targetID, targetType string) error
	Unlike(ctx context.Context, userID, targetID, targetType string) error
	IsLiked(ctx context.Context, userID, targetID, targetType string) (bool, error)
	GetLikeList(ctx context.Context, userID string, page, pageSize int32) ([]*models.Like, error)
}