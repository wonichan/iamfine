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
	
	// 获取点赞数量
	GetLikeCount(ctx context.Context, targetID, targetType string) (int64, error)
	
	// 获取点赞用户列表
	GetLikeUsers(ctx context.Context, targetID, targetType string, page, pageSize int32) ([]string, error)
}