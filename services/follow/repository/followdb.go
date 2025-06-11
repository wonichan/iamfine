package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"hupu/shared/models"
)

type followRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepository{
		db: db,
	}
}

func (fr *followRepository) Follow(ctx context.Context, followerID, followingID string) error {
	// 检查是否已经关注
	var existFollow models.Follow
	err := fr.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&existFollow).Error
	if err == nil {
		return fmt.Errorf("已经关注过了")
	}

	// 不能关注自己
	if followerID == followingID {
		return fmt.Errorf("不能关注自己")
	}

	// 创建关注关系
	newFollow := models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return fr.db.Create(&newFollow).Error
}

func (fr *followRepository) Unfollow(ctx context.Context, followerID, followingID string) error {
	return fr.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&models.Follow{}).Error
}

func (fr *followRepository) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	var existFollow models.Follow
	err := fr.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&existFollow).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (fr *followRepository) GetFollowList(ctx context.Context, userID string, page, pageSize int32) ([]string, error) {
	var follows []models.Follow
	offset := (page - 1) * pageSize
	err := fr.db.Where("follower_id = ?", userID).Offset(int(offset)).Limit(int(pageSize)).Find(&follows).Error
	if err != nil {
		return nil, err
	}

	// 转换为用户ID列表
	var followList []string
	for _, f := range follows {
		followList = append(followList, f.FollowingID)
	}

	return followList, nil
}

func (fr *followRepository) GetFollowerList(ctx context.Context, userID string, page, pageSize int32) ([]*models.Follow, error) {
	var follows []models.Follow
	offset := (page - 1) * pageSize
	err := fr.db.Where("following_id = ?", userID).Offset(int(offset)).Limit(int(pageSize)).Find(&follows).Error
	if err != nil {
		return nil, err
	}

	// 转换为指针切片
	var followerList []*models.Follow
	for i := range follows {
		followerList = append(followerList, &follows[i])
	}

	return followerList, nil
}