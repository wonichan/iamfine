package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"

	"hupu/shared/models"
	"hupu/shared/utils"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: utils.GetDB(),
	}
}
func (ur *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// 检查用户名是否已存在
	var existUser models.User
	err := ur.db.Where("username = ? ", user.Username).First(&existUser).Error
	if err == nil {
		return nil, fmt.Errorf("username %s already exists", user.Username)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// 生成用户ID
	userID := xid.New().String()

	// 密码加密
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	// 创建用户
	newUser := &models.User{
		ID:        userID,
		Username:  user.Username,
		Password:  hashedPassword,
		Nickname:  user.Username,
		Phone:     user.Phone,
		Email:     user.Email,
		Status:    models.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = ur.db.Create(newUser).Error
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (ur *UserRepository) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	// 查找用户
	var userModel models.User
	if user.ID != "" {
		if err := ur.db.Where("id = ? and status = ? ", user.ID, user.Status).First(&userModel).Error; err == nil {
			return &userModel, nil
		}
	}

	if err := ur.db.Where("username = ? and status = ? ", user.Username, user.Status).First(&userModel).Error; err != nil {
		return nil, err
	}
	return &userModel, nil
}

func (ur *UserRepository) GetUserByUsername(ctx context.Context, user *models.User) (*models.User, error) {
	var userModel models.User
	if err := ur.db.Where("id = ? ", user.ID).First(&userModel).Error; err != nil {
		return nil, err
	}
	return &userModel, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := ur.db.Model(user).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) DeleteUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := ur.db.Where("id = ?", user.ID).Delete(&models.User{}).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// 关注功能相关方法
func (ur *UserRepository) FollowUser(ctx context.Context, userID, targetUserID string) error {
	// 防止自己关注自己
	if userID == targetUserID {
		return fmt.Errorf("cannot follow yourself")
	}

	// 检查目标用户是否存在
	var targetUser models.User
	err := ur.db.Where("id = ? AND status = ?", targetUserID, models.UserStatusActive).First(&targetUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("target user not found")
		}
		return fmt.Errorf("database error when checking target user: %w", err)
	}

	// 检查是否已经关注
	var existFollow models.Follow
	err = ur.db.Where("follower_id = ? AND following_id = ?", userID, targetUserID).First(&existFollow).Error
	if err == nil {
		return fmt.Errorf("already following this user")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("database error when checking follow relationship: %w", err)
	}

	// 使用事务创建关注关系
	tx := ur.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建关注关系
	newFollow := models.Follow{
		ID:          xid.New().String(),
		FollowerID:  userID,
		FollowingID: targetUserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := tx.Create(&newFollow).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create follow relationship: %w", err)
	}

	return tx.Commit().Error
}

func (ur *UserRepository) UnfollowUser(ctx context.Context, userID, targetUserID string) error {
	err := ur.db.Where("follower_id = ? AND following_id = ?", userID, targetUserID).Delete(&models.Follow{}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func (ur *UserRepository) GetFollowerList(ctx context.Context, userID string, page, pageSize int32) ([]*models.User, error) {
	var users []*models.User
	offset := (page - 1) * pageSize

	// 使用JOIN查询优化，避免N+1查询问题
	err := ur.db.Table("users").
		Select("users.*").
		Joins("INNER JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ? AND users.status = ?", userID, models.UserStatusActive).
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) GetFollowingList(ctx context.Context, userID string, page, pageSize int32) ([]*models.User, error) {
	var users []*models.User
	offset := (page - 1) * pageSize

	// 使用JOIN查询优化，避免N+1查询问题
	err := ur.db.Table("users").
		Select("users.*").
		Joins("INNER JOIN follows ON users.id = follows.following_id").
		Where("follows.follower_id = ? AND users.status = ?", userID, models.UserStatusActive).
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// 匿名马甲管理相关方法
func (ur *UserRepository) CreateAnonymousAvatar(ctx context.Context, avatar *models.AnonymousAvatar) error {
	// 生成头像ID
	if avatar.ID == "" {
		avatar.ID = xid.New().String()
	}
	return ur.db.Create(avatar).Error
}

func (ur *UserRepository) GetAnonymousAvatarList(ctx context.Context, userID string) ([]*models.AnonymousAvatar, error) {
	var avatars []models.AnonymousAvatar
	err := ur.db.Where("user_id = ?", userID).Find(&avatars).Error
	if err != nil {
		return nil, err
	}

	// 转换为指针切片
	var avatarList []*models.AnonymousAvatar
	for i := range avatars {
		avatarList = append(avatarList, &avatars[i])
	}

	return avatarList, nil
}

func (ur *UserRepository) UpdateAnonymousAvatar(ctx context.Context, avatar *models.AnonymousAvatar) error {
	return ur.db.Model(avatar).Where("id = ?", avatar.ID).Updates(avatar).Error
}

func (ur *UserRepository) GetAnonymousAvatar(ctx context.Context, avatarID string) (*models.AnonymousAvatar, error) {
	var avatar models.AnonymousAvatar
	err := ur.db.Where("id = ?", avatarID).First(&avatar).Error
	if err != nil {
		return nil, err
	}
	return &avatar, nil
}

// 用户统计相关方法
func (ur *UserRepository) GetUserStats(ctx context.Context, userID string) (*models.UserStats, error) {
	var stats models.UserStats
	err := ur.db.Where("user_id = ?", userID).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}
