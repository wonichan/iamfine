package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hupu/shared/models"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
)

type postRedisRepo struct {
	rdb *redis.Client
}

// NewPostRedisRepo 创建新的帖子Redis仓库
func NewPostRedisRepo(rdb *redis.Client) PostRepository {
	return &postRedisRepo{
		rdb: rdb,
	}
}

// CreatePost 创建帖子
func (r *postRedisRepo) CreatePost(ctx context.Context, post *models.Post) error {
	// 生成帖子ID
	post.ID = xid.New().String()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	// 保存帖子到Redis
	err := r.savePost(ctx, post)
	if err != nil {
		return err
	}

	// 添加到用户的帖子列表
	userPostsKey := fmt.Sprintf("user:posts:%s", post.UserID)
	r.rdb.LPush(ctx, userPostsKey, post.ID)
	r.rdb.Expire(ctx, userPostsKey, 24*time.Hour)

	// 添加到全局帖子列表
	globalPostsKey := "posts:list"
	r.rdb.LPush(ctx, globalPostsKey, post.ID)
	r.rdb.Expire(ctx, globalPostsKey, 24*time.Hour)

	return nil
}

// GetPost 根据ID获取帖子
func (r *postRedisRepo) GetPost(ctx context.Context, id string) (*models.Post, error) {
	key := fmt.Sprintf("post:%s", id)
	data, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("帖子不存在")
		}
		return nil, err
	}

	var post models.Post
	err = json.Unmarshal([]byte(data), &post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

// GetPostList 获取帖子列表
func (r *postRedisRepo) GetPostList(ctx context.Context, page, pageSize int64) ([]*models.Post, error) {
	key := "posts:list"
	start := (page - 1) * pageSize
	end := start + pageSize - 1

	postIDs, err := r.rdb.LRange(ctx, key, int64(start), int64(end)).Result()
	if err != nil {
		return nil, err
	}

	var posts []*models.Post
	for _, idStr := range postIDs {
		post, err := r.GetPost(ctx, idStr)
		if err != nil {
			continue
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// UpdatePost 更新帖子
func (r *postRedisRepo) UpdatePost(ctx context.Context, post *models.Post) error {
	post.UpdatedAt = time.Now()
	return r.savePost(ctx, post)
}

// DeletePost 删除帖子
func (r *postRedisRepo) DeletePost(ctx context.Context, id string) error {
	key := fmt.Sprintf("post:%s", id)

	// 先获取帖子信息以便从列表中删除
	post, err := r.GetPost(ctx, id)
	if err != nil {
		return err
	}

	// 使用Pipeline删除相关数据
	pipe := r.rdb.Pipeline()

	// 删除帖子详情
	pipe.Del(ctx, key)

	// 从用户帖子列表中删除
	userPostsKey := fmt.Sprintf("user:posts:%s", post.UserID)
	pipe.LRem(ctx, userPostsKey, 0, id)

	// 从全局帖子列表中删除
	globalPostsKey := "posts:list"
	pipe.LRem(ctx, globalPostsKey, 0, id)

	_, err = pipe.Exec(ctx)
	return err
}

// GetPostsByUserID 根据用户ID获取帖子列表
func (r *postRedisRepo) GetPostsByUserID(ctx context.Context, userID string, page, pageSize int64) ([]*models.Post, error) {
	key := fmt.Sprintf("user:posts:%s", userID)
	start := (page - 1) * pageSize
	end := start + pageSize - 1

	postIDs, err := r.rdb.LRange(ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}

	var posts []*models.Post
	for _, idStr := range postIDs {
		post, err := r.GetPost(ctx, idStr)
		if err != nil {
			continue
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// savePost 保存帖子到Redis
func (r *postRedisRepo) savePost(ctx context.Context, post *models.Post) error {
	key := fmt.Sprintf("post:%s", post.ID)
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}

	return r.rdb.Set(ctx, key, data, 24*time.Hour).Err()
}
