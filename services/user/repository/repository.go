package repository

import (
	"context"

	"hupu/shared/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, models *models.User) (*models.User, error)
	GetUser(ctx context.Context, models *models.User) (*models.User, error)
	GetUserByUsername(ctx context.Context, models *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, models *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, models *models.User) (*models.User, error)
	
	// 关注功能相关方法
	FollowUser(ctx context.Context, userID, targetUserID string) error
	UnfollowUser(ctx context.Context, userID, targetUserID string) error
	GetFollowerList(ctx context.Context, userID string, page, pageSize int32) ([]*models.User, error)
	GetFollowingList(ctx context.Context, userID string, page, pageSize int32) ([]*models.User, error)
	
	// 匿名马甲管理相关方法
	CreateAnonymousAvatar(ctx context.Context, avatar *models.AnonymousAvatar) error
	GetAnonymousAvatarList(ctx context.Context, userID string) ([]*models.AnonymousAvatar, error)
	UpdateAnonymousAvatar(ctx context.Context, avatar *models.AnonymousAvatar) error
	GetAnonymousAvatar(ctx context.Context, avatarID string) (*models.AnonymousAvatar, error)
	
	// 用户统计相关方法
	GetUserStats(ctx context.Context, userID string) (*models.UserStats, error)
}
