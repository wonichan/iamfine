package repository

import (
	"context"

	"hupu/shared/models"
)

type FollowRepository interface {
	Follow(ctx context.Context, followerID, followingID string) error
	Unfollow(ctx context.Context, followerID, followingID string) error
	IsFollowing(ctx context.Context, followerID, followingID string) (bool, error)
	GetFollowList(ctx context.Context, userID string, page, pageSize int32) ([]string, error)
	GetFollowerList(ctx context.Context, userID string, page, pageSize int32) ([]*models.Follow, error)

	// 获取关注数量
	GetFollowCount(ctx context.Context, userID string) (int64, error)

	// 获取粉丝数量
	GetFollowerCount(ctx context.Context, userID string) (int64, error)

	// 获取共同关注
	GetMutualFollows(ctx context.Context, userID, targetUserID string, page, pageSize int32) ([]string, error)
}
