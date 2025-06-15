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

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
func (ur *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// 检查用户名是否已存在
	var existUser models.User
	err := ur.db.Where("username = ? OR phone = ?", user.Username, user.Phone).First(&existUser).Error
	if err == nil {
		return nil, fmt.Errorf("用户名或手机号已存在")
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

func (ur *userRepository) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	// 查找用户
	var userModel models.User
	err := ur.db.Where("username = ? and status = ? ", user.Username, user.Status).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}
	// 验证密码
	if !utils.CheckPassword(user.Password, userModel.Password) {
		return nil, fmt.Errorf("密码错误")
	}

	return &userModel, nil
}

func (ur *userRepository) GetUserByUsername(ctx context.Context, user *models.User) (*models.User, error) {
	var userModel models.User
	err := ur.db.Where("id = ? ", user.ID).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}
	return &userModel, nil
}

func (ur *userRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	err := ur.db.Model(user).Where("id = ?", user.ID).Updates(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) DeleteUser(ctx context.Context, user *models.User) (*models.User, error) {
	err := ur.db.Where("id = ?", user.ID).Delete(&models.User{}).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// 关注功能相关方法
func (ur *userRepository) FollowUser(ctx context.Context, userID, targetUserID string) error {
	// 检查是否已经关注
	var existFollow models.Follow
	err := ur.db.Where("follower_id = ? AND following_id = ?", userID, targetUserID).First(&existFollow).Error
	if err == nil {
		return fmt.Errorf("已经关注过了")
	}

	// 不能关注自己
	if userID == targetUserID {
		return fmt.Errorf("不能关注自己")
	}

	// 创建关注关系
	newFollow := models.Follow{
		ID:          xid.New().String(),
		FollowerID:  userID,
		FollowingID: targetUserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return ur.db.Create(&newFollow).Error
}

func (ur *userRepository) UnfollowUser(ctx context.Context, userID, targetUserID string) error {
	return ur.db.Where("follower_id = ? AND following_id = ?", userID, targetUserID).Delete(&models.Follow{}).Error
}

func (ur *userRepository) GetFollowerList(ctx context.Context, userID string, page, pageSize int32) ([]*models.User, error) {
	var follows []models.Follow
	offset := (page - 1) * pageSize
	err := ur.db.Where("following_id = ?", userID).Offset(int(offset)).Limit(int(pageSize)).Find(&follows).Error
	if err != nil {
		return nil, err
	}

	// 获取粉丝用户信息
	var users []*models.User
	for _, follow := range follows {
		var user models.User
		err := ur.db.Where("id = ?", follow.FollowerID).First(&user).Error
		if err == nil {
			users = append(users, &user)
		}
	}

	return users, nil
}

func (ur *userRepository) GetFollowingList(ctx context.Context, userID string, page, pageSize int32) ([]*models.User, error) {
	var follows []models.Follow
	offset := (page - 1) * pageSize
	err := ur.db.Where("follower_id = ?", userID).Offset(int(offset)).Limit(int(pageSize)).Find(&follows).Error
	if err != nil {
		return nil, err
	}

	// 获取关注用户信息
	var users []*models.User
	for _, follow := range follows {
		var user models.User
		err := ur.db.Where("id = ?", follow.FollowingID).First(&user).Error
		if err == nil {
			users = append(users, &user)
		}
	}

	return users, nil
}

// 匿名马甲管理相关方法
func (ur *userRepository) CreateAnonymousAvatar(ctx context.Context, avatar *models.AnonymousAvatar) error {
	// 生成头像ID
	if avatar.ID == "" {
		avatar.ID = xid.New().String()
	}
	return ur.db.Create(avatar).Error
}

func (ur *userRepository) GetAnonymousAvatarList(ctx context.Context, userID string) ([]*models.AnonymousAvatar, error) {
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

func (ur *userRepository) UpdateAnonymousAvatar(ctx context.Context, avatar *models.AnonymousAvatar) error {
	return ur.db.Model(avatar).Where("id = ?", avatar.ID).Updates(avatar).Error
}

func (ur *userRepository) GetAnonymousAvatar(ctx context.Context, avatarID string) (*models.AnonymousAvatar, error) {
	var avatar models.AnonymousAvatar
	err := ur.db.Where("id = ?", avatarID).First(&avatar).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("匿名头像不存在")
		}
		return nil, err
	}
	return &avatar, nil
}

// 用户统计相关方法
func (ur *userRepository) GetUserStats(ctx context.Context, userID string) (*models.UserStats, error) {
	var stats models.UserStats
	err := ur.db.Where("user_id = ?", userID).First(&stats).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户统计不存在")
		}
		return nil, err
	}
	return &stats, nil
}
