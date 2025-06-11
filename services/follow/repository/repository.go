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
}