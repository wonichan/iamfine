package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"

	"hupu/shared/models"
	"hupu/shared/utils"
)

type userRedisRepo struct {
	rdb *redis.Client
}

func NewUserRedisRepo(rdb *redis.Client) UserRepository {
	return &userRedisRepo{
		rdb: rdb,
	}
}

func (ur *userRedisRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// 检查用户名是否已存在
	existUser, err := ur.getUserByField(ctx, "username", user.Username)
	if err == nil && existUser != nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 检查手机号是否已存在
	if user.Phone != "" {
		existUser, err = ur.getUserByField(ctx, "phone", user.Phone)
		if err == nil && existUser != nil {
			return nil, fmt.Errorf("手机号已存在")
		}
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

	// 存储到Redis
	err = ur.saveUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (ur *userRedisRepo) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	// 根据用户名查找用户
	userModel, err := ur.getUserByField(ctx, "username", user.Username)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if userModel == nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 验证密码
	if !utils.CheckPassword(user.Password, userModel.Password) {
		return nil, fmt.Errorf("密码错误")
	}

	return userModel, nil
}

func (ur *userRedisRepo) GetUserByUsername(ctx context.Context, user *models.User) (*models.User, error) {
	// 根据用户ID查找用户
	userModel, err := ur.getUserByField(ctx, "id", user.ID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if userModel == nil {
		return nil, fmt.Errorf("用户不存在")
	}

	return userModel, nil
}

func (ur *userRedisRepo) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// 检查用户是否存在
	existUser, err := ur.getUserByField(ctx, "id", user.ID)
	if err != nil || existUser == nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 更新字段
	if user.Nickname != "" {
		existUser.Nickname = user.Nickname
	}
	if user.Avatar != "" {
		existUser.Avatar = user.Avatar
	}
	if user.Email != "" {
		existUser.Email = user.Email
	}
	if user.Phone != "" {
		existUser.Phone = user.Phone
	}
	existUser.UpdatedAt = time.Now()

	// 保存更新后的用户
	err = ur.saveUser(ctx, existUser)
	if err != nil {
		return nil, err
	}

	return existUser, nil
}

func (ur *userRedisRepo) DeleteUser(ctx context.Context, user *models.User) (*models.User, error) {
	// 检查用户是否存在
	existUser, err := ur.getUserByField(ctx, "id", user.ID)
	if err != nil || existUser == nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 删除用户数据
	userKey := fmt.Sprintf("user:%s", user.ID)
	usernameKey := fmt.Sprintf("username:%s", existUser.Username)
	phoneKey := fmt.Sprintf("phone:%s", existUser.Phone)

	pipe := ur.rdb.Pipeline()
	pipe.Del(ctx, userKey)
	pipe.Del(ctx, usernameKey)
	if existUser.Phone != "" {
		pipe.Del(ctx, phoneKey)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return existUser, nil
}

// 辅助方法：根据字段查找用户
func (ur *userRedisRepo) getUserByField(ctx context.Context, field, value string) (*models.User, error) {
	var userID string
	var err error

	switch field {
	case "id":
		userID = value
	case "username":
		usernameKey := fmt.Sprintf("username:%s", value)
		userID, err = ur.rdb.Get(ctx, usernameKey).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil, fmt.Errorf("用户不存在")
			}
			return nil, err
		}
	case "phone":
		phoneKey := fmt.Sprintf("phone:%s", value)
		userID, err = ur.rdb.Get(ctx, phoneKey).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil, fmt.Errorf("用户不存在")
			}
			return nil, err
		}
	default:
		return nil, fmt.Errorf("不支持的查询字段")
	}

	// 根据用户ID获取用户信息
	userKey := fmt.Sprintf("user:%s", userID)
	userData, err := ur.rdb.Get(ctx, userKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}

	// 反序列化用户数据
	var user models.User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// 辅助方法：保存用户到Redis
func (ur *userRedisRepo) saveUser(ctx context.Context, user *models.User) error {
	// 序列化用户数据
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// 使用Pipeline批量操作
	pipe := ur.rdb.Pipeline()

	// 存储用户数据
	userKey := fmt.Sprintf("user:%s", user.ID)
	pipe.Set(ctx, userKey, userData, 24*time.Hour) // 设置24小时过期时间

	// 建立用户名到用户ID的映射
	usernameKey := fmt.Sprintf("username:%s", user.Username)
	pipe.Set(ctx, usernameKey, user.ID, 24*time.Hour)

	// 建立手机号到用户ID的映射（如果有手机号）
	if user.Phone != "" {
		phoneKey := fmt.Sprintf("phone:%s", user.Phone)
		pipe.Set(ctx, phoneKey, user.ID, 24*time.Hour)
	}

	// 执行Pipeline
	_, err = pipe.Exec(ctx)
	return err
}
