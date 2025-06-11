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
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	// 检查手机号是否已存在
	err = ur.db.Where("phone = ?", user.Phone).First(&existUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return newUser, ur.db.Create(user).Error
}

func (ur *userRepository) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	// 查找用户
	var userModel models.User
	err := ur.db.Where("username = ?", user.Username).First(&userModel).Error
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
	panic("not implemented") // TODO: Implement
}

func (ur *userRepository) DeleteUser(ctx context.Context, models *models.User) (*models.User, error) {
	panic("not implemented") // TODO: Implement
}
