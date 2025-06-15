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
		Status:    models.UserStatusActive,
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
	var (
		userModel *models.User
		err       error
	)
	if user.ID != "" {
		// 根据用户ID查找用户
		userModel, err = ur.getUserByField(ctx, "id", user.ID)
	} else {
		// 根据用户名查找用户
		userModel, err = ur.getUserByField(ctx, "username", user.Username)
	}
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if userModel == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if user.Status != models.UserStatusActive {
		if user.Status != userModel.Status {
			return nil, fmt.Errorf("用户状态异常")
		}
	}
	// 验证密码
	if !utils.CheckPassword(user.Password, userModel.Password) {
		return nil, fmt.Errorf("密码错误")
	}

	return userModel, nil
}

func (ur *userRedisRepo) GetUserByUsername(ctx context.Context, user *models.User) (*models.User, error) {
	// 根据用户ID查找用户
	userModel, err := ur.getUserByField(ctx, "username", user.Username)
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

// 关注功能相关方法
func (ur *userRedisRepo) FollowUser(ctx context.Context, userID, targetUserID string) error {
	// 检查是否已经关注
	followKey := fmt.Sprintf("follow:%s:%s", userID, targetUserID)
	exists, err := ur.rdb.Exists(ctx, followKey).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return fmt.Errorf("已经关注过了")
	}

	// 不能关注自己
	if userID == targetUserID {
		return fmt.Errorf("不能关注自己")
	}

	// 创建关注关系
	follow := &models.Follow{
		ID:          xid.New().String(),
		FollowerID:  userID,
		FollowingID: targetUserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存关注关系
	followData, err := json.Marshal(follow)
	if err != nil {
		return err
	}

	followingListKey := fmt.Sprintf("follow:list:%s", userID)
	followerListKey := fmt.Sprintf("follower:list:%s", targetUserID)

	pipe := ur.rdb.Pipeline()
	pipe.Set(ctx, followKey, followData, 24*time.Hour)
	pipe.LPush(ctx, followingListKey, targetUserID)
	pipe.LPush(ctx, followerListKey, userID)
	pipe.Expire(ctx, followingListKey, 24*time.Hour)
	pipe.Expire(ctx, followerListKey, 24*time.Hour)
	_, err = pipe.Exec(ctx)

	return err
}

func (ur *userRedisRepo) UnfollowUser(ctx context.Context, userID, targetUserID string) error {
	followKey := fmt.Sprintf("follow:%s:%s", userID, targetUserID)
	followingListKey := fmt.Sprintf("follow:list:%s", userID)
	followerListKey := fmt.Sprintf("follower:list:%s", targetUserID)

	pipe := ur.rdb.Pipeline()
	pipe.Del(ctx, followKey)
	pipe.LRem(ctx, followingListKey, 0, targetUserID)
	pipe.LRem(ctx, followerListKey, 0, userID)
	_, err := pipe.Exec(ctx)

	return err
}

func (ur *userRedisRepo) GetFollowerList(ctx context.Context, userID string, page, pageSize int32) ([]*models.User, error) {
	followerListKey := fmt.Sprintf("follower:list:%s", userID)

	// 分页获取粉丝列表
	start := (page - 1) * pageSize
	end := start + pageSize - 1
	followerIDs, err := ur.rdb.LRange(ctx, followerListKey, int64(start), int64(end)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []*models.User{}, nil
		}
		return nil, err
	}

	// 批量获取用户信息
	var users []*models.User
	for _, followerID := range followerIDs {
		userKey := fmt.Sprintf("user:%s", followerID)
		userData, err := ur.rdb.Get(ctx, userKey).Result()
		if err == nil {
			var user models.User
			if json.Unmarshal([]byte(userData), &user) == nil {
				users = append(users, &user)
			}
		}
	}

	return users, nil
}

func (ur *userRedisRepo) GetFollowingList(ctx context.Context, userID string, page, pageSize int32) ([]*models.User, error) {
	followingListKey := fmt.Sprintf("follow:list:%s", userID)

	// 分页获取关注列表
	start := (page - 1) * pageSize
	end := start + pageSize - 1
	followingIDs, err := ur.rdb.LRange(ctx, followingListKey, int64(start), int64(end)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []*models.User{}, nil
		}
		return nil, err
	}

	// 批量获取用户信息
	var users []*models.User
	for _, followingID := range followingIDs {
		userKey := fmt.Sprintf("user:%s", followingID)
		userData, err := ur.rdb.Get(ctx, userKey).Result()
		if err == nil {
			var user models.User
			if json.Unmarshal([]byte(userData), &user) == nil {
				users = append(users, &user)
			}
		}
	}

	return users, nil
}

// 匿名马甲管理相关方法
func (ur *userRedisRepo) CreateAnonymousAvatar(ctx context.Context, avatar *models.AnonymousAvatar) error {
	// 生成头像ID
	if avatar.ID == "" {
		avatar.ID = xid.New().String()
	}
	
	avatarData, err := json.Marshal(avatar)
	if err != nil {
		return err
	}

	avatarKey := fmt.Sprintf("anonymous_avatar:%s", avatar.ID)
	avatarListKey := fmt.Sprintf("anonymous_avatar:list:%s", avatar.UserID)

	pipe := ur.rdb.Pipeline()
	pipe.Set(ctx, avatarKey, avatarData, 24*time.Hour)
	pipe.LPush(ctx, avatarListKey, avatar.ID)
	pipe.Expire(ctx, avatarListKey, 24*time.Hour)
	_, err = pipe.Exec(ctx)

	return err
}

func (ur *userRedisRepo) GetAnonymousAvatarList(ctx context.Context, userID string) ([]*models.AnonymousAvatar, error) {
	avatarListKey := fmt.Sprintf("anonymous_avatar:list:%s", userID)
	avatarIDs, err := ur.rdb.LRange(ctx, avatarListKey, 0, -1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []*models.AnonymousAvatar{}, nil
		}
		return nil, err
	}

	// 批量获取头像信息
	var avatars []*models.AnonymousAvatar
	for _, avatarID := range avatarIDs {
		avatarKey := fmt.Sprintf("anonymous_avatar:%s", avatarID)
		avatarData, err := ur.rdb.Get(ctx, avatarKey).Result()
		if err == nil {
			var avatar models.AnonymousAvatar
			if json.Unmarshal([]byte(avatarData), &avatar) == nil {
				avatars = append(avatars, &avatar)
			}
		}
	}

	return avatars, nil
}

func (ur *userRedisRepo) UpdateAnonymousAvatar(ctx context.Context, avatar *models.AnonymousAvatar) error {
	avatarData, err := json.Marshal(avatar)
	if err != nil {
		return err
	}

	avatarKey := fmt.Sprintf("anonymous_avatar:%s", avatar.ID)
	return ur.rdb.Set(ctx, avatarKey, avatarData, 24*time.Hour).Err()
}

func (ur *userRedisRepo) GetAnonymousAvatar(ctx context.Context, avatarID string) (*models.AnonymousAvatar, error) {
	avatarKey := fmt.Sprintf("anonymous_avatar:%s", avatarID)
	avatarData, err := ur.rdb.Get(ctx, avatarKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("匿名头像不存在")
		}
		return nil, err
	}

	// 反序列化头像数据
	var avatar models.AnonymousAvatar
	err = json.Unmarshal([]byte(avatarData), &avatar)
	if err != nil {
		return nil, err
	}

	return &avatar, nil
}

// 用户统计相关方法
func (ur *userRedisRepo) GetUserStats(ctx context.Context, userID string) (*models.UserStats, error) {
	statsKey := fmt.Sprintf("user_stats:%s", userID)
	statsData, err := ur.rdb.Get(ctx, statsKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("用户统计不存在")
		}
		return nil, err
	}

	// 反序列化统计数据
	var stats models.UserStats
	err = json.Unmarshal([]byte(statsData), &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
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
