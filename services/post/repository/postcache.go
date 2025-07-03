package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hupu/shared/models"
	"time"

	"github.com/redis/go-redis/v9"
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

// GetPostListByTopic 根据话题ID获取帖子列表
func (r *postRedisRepo) GetPostListByTopic(ctx context.Context, topicID string, page, pageSize int64) ([]*models.Post, error) {
	key := fmt.Sprintf("topic:posts:%s", topicID)
	start := (page - 1) * pageSize
	end := start + pageSize - 1

	postIDs, err := r.rdb.LRange(ctx, key, start, end).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []*models.Post{}, nil
		}
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

// GetPostListByCategory 根据分类获取帖子列表
func (r *postRedisRepo) GetPostListByCategory(ctx context.Context, category string, page, pageSize int64) ([]*models.Post, error) {
	key := fmt.Sprintf("category:posts:%s", category)
	start := (page - 1) * pageSize
	end := start + pageSize - 1

	postIDs, err := r.rdb.LRange(ctx, key, start, end).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []*models.Post{}, nil
		}
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

// GetPostListWithConditions 根据多个条件获取帖子列表 (Redis实现暂时回退到数据库查询)
func (r *postRedisRepo) GetPostListWithConditions(ctx context.Context, conditions map[string]interface{}, page, pageSize int64, sortType string) ([]*models.Post, error) {
	// Redis缓存的复合查询比较复杂，这里暂时返回空结果
	// 在实际使用中，handler会回退到数据库查询
	return []*models.Post{}, nil
}

// IncrementViewCount 增加浏览次数
func (r *postRedisRepo) IncrementViewCount(ctx context.Context, postID string) error {
	// 获取当前帖子数据
	post, err := r.GetPost(ctx, postID)
	if err != nil {
		return err
	}
	
	// 增加浏览次数
	post.ViewCount++
	post.UpdatedAt = time.Now()
	
	// 保存更新后的帖子
	return r.savePost(ctx, post)
}

// 话题管理相关方法
func (r *postRedisRepo) CreateTopic(ctx context.Context, topic *models.Topic) error {
	topic.ID = xid.New().String()
	topic.CreatedAt = time.Now()
	topic.UpdatedAt = time.Now()

	// 保存话题到Redis
	err := r.saveTopic(ctx, topic)
	if err != nil {
		return err
	}

	// 添加到话题列表
	topicsKey := "topics:list"
	r.rdb.LPush(ctx, topicsKey, topic.ID)
	r.rdb.Expire(ctx, topicsKey, 24*time.Hour)

	return nil
}

func (r *postRedisRepo) GetTopicList(ctx context.Context, page, pageSize int32) ([]*models.Topic, error) {
	key := "topics:list"
	start := (page - 1) * pageSize
	end := start + pageSize - 1

	topicIDs, err := r.rdb.LRange(ctx, key, int64(start), int64(end)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []*models.Topic{}, nil
		}
		return nil, err
	}

	var topics []*models.Topic
	for _, idStr := range topicIDs {
		topic, err := r.GetTopic(ctx, idStr)
		if err != nil {
			continue
		}
		topics = append(topics, topic)
	}

	return topics, nil
}

func (r *postRedisRepo) GetTopic(ctx context.Context, topicID string) (*models.Topic, error) {
	key := fmt.Sprintf("topic:%s", topicID)
	data, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("话题不存在")
		}
		return nil, err
	}

	var topic models.Topic
	err = json.Unmarshal([]byte(data), &topic)
	if err != nil {
		return nil, err
	}

	return &topic, nil
}

// 收藏功能相关方法
func (r *postRedisRepo) FavoritePost(ctx context.Context, userID string, postID string) error {
	// 检查是否已经收藏
	favoriteKey := fmt.Sprintf("favorite:%s:%s", userID, postID)
	exists, err := r.rdb.Exists(ctx, favoriteKey).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return nil // 已经收藏过了
	}

	// 创建收藏记录
	favorite := &models.PostFavorite{
		ID:        xid.New().String(),
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now(),
	}

	// 保存收藏记录
	favoriteData, err := json.Marshal(favorite)
	if err != nil {
		return err
	}

	favoriteListKey := fmt.Sprintf("favorite:list:%s", userID)

	pipe := r.rdb.Pipeline()
	pipe.Set(ctx, favoriteKey, favoriteData, 24*time.Hour)
	pipe.LPush(ctx, favoriteListKey, postID)
	pipe.Expire(ctx, favoriteListKey, 24*time.Hour)
	_, err = pipe.Exec(ctx)

	return err
}

func (r *postRedisRepo) UnfavoritePost(ctx context.Context, userID string, postID string) error {
	favoriteKey := fmt.Sprintf("favorite:%s:%s", userID, postID)
	favoriteListKey := fmt.Sprintf("favorite:list:%s", userID)

	pipe := r.rdb.Pipeline()
	pipe.Del(ctx, favoriteKey)
	pipe.LRem(ctx, favoriteListKey, 0, postID)
	_, err := pipe.Exec(ctx)

	return err
}

func (r *postRedisRepo) GetFavoriteList(ctx context.Context, userID string, page, pageSize int32) ([]*models.Post, error) {
	key := fmt.Sprintf("favorite:list:%s", userID)
	start := (page - 1) * pageSize
	end := start + pageSize - 1

	postIDs, err := r.rdb.LRange(ctx, key, int64(start), int64(end)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []*models.Post{}, nil
		}
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

// 评分功能相关方法
func (r *postRedisRepo) RatePost(ctx context.Context, rating *models.PostRating) error {
	// 检查是否已经评分过
	ratingKey := fmt.Sprintf("rating:%s:%s", rating.UserID, rating.PostID)
	existingData, err := r.rdb.Get(ctx, ratingKey).Result()

	if err == redis.Nil {
		// 创建新评分
		rating.ID = xid.New().String()
		rating.CreatedAt = time.Now()
	} else if err != nil {
		return err
	} else {
		// 更新现有评分
		var existingRating models.PostRating
		if json.Unmarshal([]byte(existingData), &existingRating) == nil {
			rating.ID = existingRating.ID
			rating.CreatedAt = existingRating.CreatedAt
		}
		rating.UpdatedAt = time.Now()
	}

	// 保存评分
	ratingData, err := json.Marshal(rating)
	if err != nil {
		return err
	}

	return r.rdb.Set(ctx, ratingKey, ratingData, 24*time.Hour).Err()
}

func (r *postRedisRepo) GetScoreRanking(ctx context.Context, page, pageSize int32) ([]*models.Post, error) {
	// Redis中的评分排名实现相对复杂，这里简化处理
	// 实际项目中可能需要使用有序集合(ZSET)来维护排名
	key := "posts:ranking"
	start := (page - 1) * pageSize
	end := start + pageSize - 1

	postIDs, err := r.rdb.LRange(ctx, key, int64(start), int64(end)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []*models.Post{}, nil
		}
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

// 获取匿名头像信息
func (r *postRedisRepo) GetAnonymousAvatar(ctx context.Context, avatarID string) (*models.AnonymousAvatar, error) {
	key := fmt.Sprintf("anonymous_avatar:%s", avatarID)
	data, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("匿名头像不存在")
		}
		return nil, err
	}

	var avatar models.AnonymousAvatar
	err = json.Unmarshal([]byte(data), &avatar)
	if err != nil {
		return nil, err
	}

	return &avatar, nil
}

// saveTopic 保存话题到Redis
func (r *postRedisRepo) saveTopic(ctx context.Context, topic *models.Topic) error {
	key := fmt.Sprintf("topic:%s", topic.ID)
	data, err := json.Marshal(topic)
	if err != nil {
		return err
	}

	return r.rdb.Set(ctx, key, data, 24*time.Hour).Err()
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
