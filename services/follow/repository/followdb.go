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

func (fr *followRepository) GetFollowCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := fr.db.Model(&models.Follow{}).Where("follower_id = ?", userID).Count(&count).Error
	return count, err
}

func (fr *followRepository) GetFollowerCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := fr.db.Model(&models.Follow{}).Where("following_id = ?", userID).Count(&count).Error
	return count, err
}

func (fr *followRepository) GetMutualFollows(ctx context.Context, userID1, userID2 string, page, pageSize int32) ([]string, error) {
	// 查找用户1关注的人中，也关注用户2的人
	var mutualFollows []string
	offset := (page - 1) * pageSize

	// 子查询：用户1关注的人
	subQuery1 := fr.db.Model(&models.Follow{}).Select("following_id").Where("follower_id = ?", userID1)
	// 子查询：用户2关注的人
	subQuery2 := fr.db.Model(&models.Follow{}).Select("following_id").Where("follower_id = ?", userID2)

	// 查找交集
	err := fr.db.Model(&models.Follow{}).
		Select("following_id").
		Where("follower_id = ? AND following_id IN (?)", subQuery1, subQuery2).
		Offset(int(offset)).
		Limit(int(pageSize)).
		Pluck("following_id", &mutualFollows).Error

	return mutualFollows, err
}
