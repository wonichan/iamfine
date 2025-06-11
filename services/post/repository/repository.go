package repository

import (
	"context"
	"hupu/shared/models"
)

// PostRepository 定义帖子数据访问接口
type PostRepository interface {
	// CreatePost 创建帖子
	CreatePost(ctx context.Context, post *models.Post) error

	// GetPost 根据ID获取帖子
	GetPost(ctx context.Context, id string) (*models.Post, error)

	// GetPostList 获取帖子列表
	GetPostList(ctx context.Context, page, pageSize int64) ([]*models.Post, error)

	// UpdatePost 更新帖子
	UpdatePost(ctx context.Context, post *models.Post) error

	// DeletePost 删除帖子
	DeletePost(ctx context.Context, id string) error

	// GetPostsByUserID 根据用户ID获取帖子列表
	GetPostsByUserID(ctx context.Context, userID string, page, pageSize int64) ([]*models.Post, error)
}
