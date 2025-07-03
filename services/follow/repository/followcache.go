package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"hupu/shared/models"
)

type followRedisRepo struct {
	rdb *redis.Client
}

func NewFollowRedisRepo(rdb *redis.Client) FollowRepository {
	return &followRedisRepo{
		rdb: rdb,
	}
}

func (fr *followRedisRepo) Follow(ctx context.Context, followerID, followingID string) error {
	// 检查是否已经关注
	isFollowing, err := fr.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if isFollowing {
		return fmt.Errorf("已经关注过了")
	}

	// 不能关注自己
	if followerID == followingID {
		return fmt.Errorf("不能关注自己")
	}

	// 创建关注关系
	newFollow := &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   time.Now(),
	}

	// 保存到Redis
	return fr.saveFollow(ctx, newFollow)
}

func (fr *followRedisRepo) Unfollow(ctx context.Context, followerID, followingID string) error {
	// 构建键名
	followKey := fmt.Sprintf("follow:%s:%s", followerID, followingID)
	followingListKey := fmt.Sprintf("follow:list:%s", followerID)
	followerListKey := fmt.Sprintf("follower:list:%s", followingID)

	pipe := fr.rdb.Pipeline()
	pipe.Del(ctx, followKey)
	pipe.LRem(ctx, followingListKey, 0, followingID)
	pipe.LRem(ctx, followerListKey, 0, followerID)
	_, err := pipe.Exec(ctx)

	return err
}

func (fr *followRedisRepo) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	followKey := fmt.Sprintf("follow:%s:%s", followerID, followingID)
	exists, err := fr.rdb.Exists(ctx, followKey).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (fr *followRedisRepo) GetFollowList(ctx context.Context, userID string, page, pageSize int32) ([]string, error) {
	followingListKey := fmt.Sprintf("follow:list:%s", userID)

	// 分页获取关注列表
	start := (page - 1) * pageSize
	end := start + pageSize - 1
	followingIDs, err := fr.rdb.LRange(ctx, followingListKey, int64(start), int64(end)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []string{}, nil
		}
		return nil, err
	}

	return followingIDs, nil
}

func (fr *followRedisRepo) GetFollowerList(ctx context.Context, userID string, page, pageSize int32) ([]*models.Follow, error) {
	followerListKey := fmt.Sprintf("follower:list:%s", userID)

	// 分页获取粉丝列表
	start := (page - 1) * pageSize
	end := start + pageSize - 1
	followerIDs, err := fr.rdb.LRange(ctx, followerListKey, int64(start), int64(end)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []*models.Follow{}, nil
		}
		return nil, err
	}

	// 批量获取关注详情
	var followers []*models.Follow
	for _, followerID := range followerIDs {
		followKey := fmt.Sprintf("follow:%s:%s", followerID, userID)
		followData, err := fr.rdb.Get(ctx, followKey).Result()
		if err == nil {
			var follow models.Follow
			if json.Unmarshal([]byte(followData), &follow) == nil {
				followers = append(followers, &follow)
			}
		}
	}

	return followers, nil
}

func (fr *followRedisRepo) GetFollowCount(ctx context.Context, userID string) (int64, error) {
	followingListKey := fmt.Sprintf("follow:list:%s", userID)
	count, err := fr.rdb.LLen(ctx, followingListKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func (fr *followRedisRepo) GetFollowerCount(ctx context.Context, userID string) (int64, error) {
	followerListKey := fmt.Sprintf("follower:list:%s", userID)
	count, err := fr.rdb.LLen(ctx, followerListKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func (fr *followRedisRepo) GetMutualFollows(ctx context.Context, userID1, userID2 string, page, pageSize int32) ([]string, error) {
	// 获取两个用户的关注列表
	followingListKey1 := fmt.Sprintf("follow:list:%s", userID1)
	followingListKey2 := fmt.Sprintf("follow:list:%s", userID2)

	// 使用Redis的集合操作找交集
	tempKey := fmt.Sprintf("temp:mutual:%s:%s", userID1, userID2)

	// 将列表转换为集合并求交集
	pipe := fr.rdb.Pipeline()
	pipe.Del(ctx, tempKey)

	// 获取用户1的关注列表
	followingIDs1, err := fr.rdb.LRange(ctx, followingListKey1, 0, -1).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	// 获取用户2的关注列表
	followingIDs2, err := fr.rdb.LRange(ctx, followingListKey2, 0, -1).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	// 找出交集
	followingSet2 := make(map[string]bool)
	for _, id := range followingIDs2 {
		followingSet2[id] = true
	}

	var mutualFollows []string
	for _, id := range followingIDs1 {
		if followingSet2[id] {
			mutualFollows = append(mutualFollows, id)
		}
	}

	// 分页处理
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= int32(len(mutualFollows)) {
		return []string{}, nil
	}
	if end > int32(len(mutualFollows)) {
		end = int32(len(mutualFollows))
	}

	return mutualFollows[start:end], nil
}

// 辅助方法：保存关注关系到Redis
func (fr *followRedisRepo) saveFollow(ctx context.Context, follow *models.Follow) error {
	// 序列化关注数据
	followData, err := json.Marshal(follow)
	if err != nil {
		return err
	}

	// 构建键名
	followKey := fmt.Sprintf("follow:%s:%s", follow.FollowerID, follow.FollowingID)
	followingListKey := fmt.Sprintf("follow:list:%s", follow.FollowerID)
	followerListKey := fmt.Sprintf("follower:list:%s", follow.FollowingID)

	// 使用Pipeline批量操作
	pipe := fr.rdb.Pipeline()
	pipe.Set(ctx, followKey, followData, 24*time.Hour)
	pipe.LPush(ctx, followingListKey, follow.FollowingID)
	pipe.LPush(ctx, followerListKey, follow.FollowerID)
	pipe.Expire(ctx, followingListKey, 24*time.Hour)
	pipe.Expire(ctx, followerListKey, 24*time.Hour)
	_, err = pipe.Exec(ctx)

	return err
}
