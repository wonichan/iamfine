package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"hupu/shared/models"
)

type likeRedisRepo struct {
	rdb *redis.Client
}

func NewLikeRedisRepo(rdb *redis.Client) LikeRepository {
	return &likeRedisRepo{
		rdb: rdb,
	}
}

func (lr *likeRedisRepo) Like(ctx context.Context, userID, targetID, targetType string) error {
	// 检查是否已经点赞
	isLiked, err := lr.IsLiked(ctx, userID, targetID, targetType)
	if err != nil {
		return err
	}
	if isLiked {
		return fmt.Errorf("已经点赞过了")
	}

	// 创建点赞记录
	newLike := &models.Like{
		UserID:     userID,
		TargetID:   targetID,
		TargetType: targetType,
		CreatedAt:  time.Now(),
	}

	// 保存到Redis
	err = lr.saveLike(ctx, newLike)
	if err != nil {
		return err
	}

	// 更新目标对象的点赞数
	return lr.updateLikeCount(ctx, targetID, targetType, 1)
}

func (lr *likeRedisRepo) Unlike(ctx context.Context, userID, targetID, targetType string) error {
	// 构建键名
	likeKey := fmt.Sprintf("like:%s:%s:%s", userID, targetID, targetType)
	likeListKey := fmt.Sprintf("like:list:%s", userID)

	pipe := lr.rdb.Pipeline()
	pipe.Del(ctx, likeKey)
	pipe.LRem(ctx, likeListKey, 0, fmt.Sprintf("%s:%s", targetID, targetType))
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	// 更新目标对象的点赞数
	return lr.updateLikeCount(ctx, targetID, targetType, -1)
}

func (lr *likeRedisRepo) IsLiked(ctx context.Context, userID, targetID, targetType string) (bool, error) {
	likeKey := fmt.Sprintf("like:%s:%s:%s", userID, targetID, targetType)
	exists, err := lr.rdb.Exists(ctx, likeKey).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (lr *likeRedisRepo) GetLikeList(ctx context.Context, userID string, page, pageSize int32) ([]*models.Like, error) {
	likeListKey := fmt.Sprintf("like:list:%s", userID)

	// 分页获取点赞列表
	start := (page - 1) * pageSize
	end := start + pageSize - 1
	likeItems, err := lr.rdb.LRange(ctx, likeListKey, int64(start), int64(end)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []*models.Like{}, nil
		}
		return nil, err
	}

	// 批量获取点赞详情
	var likes []*models.Like
	for _, item := range likeItems {
		// item格式: "targetID:targetType"
		likeKey := fmt.Sprintf("like:%s:%s", userID, item)
		likeData, err := lr.rdb.Get(ctx, likeKey).Result()
		if err == nil {
			var like models.Like
			if json.Unmarshal([]byte(likeData), &like) == nil {
				likes = append(likes, &like)
			}
		}
	}

	return likes, nil
}

// 辅助方法：保存点赞到Redis
func (lr *likeRedisRepo) saveLike(ctx context.Context, like *models.Like) error {
	// 序列化点赞数据
	likeData, err := json.Marshal(like)
	if err != nil {
		return err
	}

	// 构建键名
	likeKey := fmt.Sprintf("like:%s:%s:%s", like.UserID, like.TargetID, like.TargetType)
	likeListKey := fmt.Sprintf("like:list:%s", like.UserID)

	// 使用Pipeline批量操作
	pipe := lr.rdb.Pipeline()
	pipe.Set(ctx, likeKey, likeData, 24*time.Hour)
	pipe.LPush(ctx, likeListKey, fmt.Sprintf("%s:%s", like.TargetID, like.TargetType))
	pipe.Expire(ctx, likeListKey, 24*time.Hour)
	_, err = pipe.Exec(ctx)

	return err
}

// 辅助方法：更新点赞数
func (lr *likeRedisRepo) updateLikeCount(ctx context.Context, targetID, targetType string, delta int) error {
	var countKey string
	if targetType == "post" {
		countKey = fmt.Sprintf("post:like_count:%s", targetID)
	} else if targetType == "comment" {
		countKey = fmt.Sprintf("comment:like_count:%s", targetID)
	} else {
		return fmt.Errorf("不支持的目标类型: %s", targetType)
	}

	// 更新计数器
	_, err := lr.rdb.IncrBy(ctx, countKey, int64(delta)).Result()
	if err != nil {
		return err
	}

	// 设置过期时间
	lr.rdb.Expire(ctx, countKey, 24*time.Hour)
	return nil
}