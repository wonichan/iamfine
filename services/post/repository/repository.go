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

	// GetPostListByTopic 根据话题ID获取帖子列表
	GetPostListByTopic(ctx context.Context, topicID string, page, pageSize int64) ([]*models.Post, error)

	// GetPostListByCategory 根据分类获取帖子列表
	GetPostListByCategory(ctx context.Context, category string, page, pageSize int64) ([]*models.Post, error)

	// UpdatePost 更新帖子
	UpdatePost(ctx context.Context, post *models.Post) error

	// DeletePost 删除帖子
	DeletePost(ctx context.Context, id string) error

	// GetPostsByUserID 根据用户ID获取帖子列表
	GetPostsByUserID(ctx context.Context, userID string, page, pageSize int64) ([]*models.Post, error)

	// IncrementViewCount 增加浏览次数
	IncrementViewCount(ctx context.Context, postID string) error

	// 话题管理相关方法
	CreateTopic(ctx context.Context, topic *models.Topic) error
	GetTopicList(ctx context.Context, page, pageSize int32) ([]*models.Topic, error)
	GetTopic(ctx context.Context, topicID string) (*models.Topic, error)

	// 收藏功能相关方法
	FavoritePost(ctx context.Context, userID string, postID string) error
	UnfavoritePost(ctx context.Context, userID string, postID string) error
	GetFavoriteList(ctx context.Context, userID string, page, pageSize int32) ([]*models.Post, error)

	// 评分功能相关方法
	RatePost(ctx context.Context, rating *models.PostRating) error
	GetScoreRanking(ctx context.Context, page, pageSize int32) ([]*models.Post, error)

	// 获取匿名头像信息
	GetAnonymousAvatar(ctx context.Context, avatarID string) (*models.AnonymousAvatar, error)
}
