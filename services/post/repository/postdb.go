package repository

import (
	"context"
	"hupu/shared/models"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

// NewPostRepository 创建新的帖子数据库仓库
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		db: db,
	}
}

// CreatePost 创建帖子
func (r *postRepository) CreatePost(ctx context.Context, post *models.Post) error {
	// 生成帖子ID
	post.ID = xid.New().String()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).Create(post).Error
}

// GetPost 根据ID获取帖子
func (r *postRepository) GetPost(ctx context.Context, id string) (*models.Post, error) {
	var post models.Post
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostList 获取帖子列表
func (r *postRepository) GetPostList(ctx context.Context, page, pageSize int64) ([]*models.Post, error) {
	var posts []*models.Post
	offset := (page - 1) * pageSize

	err := r.db.WithContext(ctx).
		Offset(int(offset)).
		Limit(int(pageSize)).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, err
}

// UpdatePost 更新帖子
func (r *postRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	post.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(post).Error
}

// DeletePost 删除帖子
func (r *postRepository) DeletePost(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Post{}).Error
}

// GetPostsByUserID 根据用户ID获取帖子列表
func (r *postRepository) GetPostsByUserID(ctx context.Context, userID string, page, pageSize int64) ([]*models.Post, error) {
	var posts []*models.Post
	offset := (page - 1) * pageSize

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Offset(int(offset)).
		Limit(int(pageSize)).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, err
}

// GetPostListByTopic 根据话题ID获取帖子列表
func (r *postRepository) GetPostListByTopic(ctx context.Context, topicID string, page, pageSize int64) ([]*models.Post, error) {
	var posts []*models.Post
	offset := (page - 1) * pageSize

	err := r.db.WithContext(ctx).
		Where("topic_id = ?", topicID).
		Offset(int(offset)).
		Limit(int(pageSize)).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, err
}

// GetPostListByCategory 根据分类获取帖子列表
func (r *postRepository) GetPostListByCategory(ctx context.Context, category string, page, pageSize int64) ([]*models.Post, error) {
	var posts []*models.Post
	offset := (page - 1) * pageSize

	err := r.db.WithContext(ctx).
		Where("category = ?", category).
		Offset(int(offset)).
		Limit(int(pageSize)).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, err
}

// IncrementViewCount 增加浏览次数
func (r *postRepository) IncrementViewCount(ctx context.Context, postID string) error {
	return r.db.WithContext(ctx).
		Model(&models.Post{}).
		Where("id = ?", postID).
		Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// 话题管理相关方法
func (r *postRepository) CreateTopic(ctx context.Context, topic *models.Topic) error {
	topic.ID = xid.New().String()
	topic.CreatedAt = time.Now()
	topic.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(topic).Error
}

func (r *postRepository) GetTopicList(ctx context.Context, page, pageSize int32) ([]*models.Topic, error) {
	var topics []*models.Topic
	offset := (page - 1) * pageSize

	err := r.db.WithContext(ctx).
		Offset(int(offset)).
		Limit(int(pageSize)).
		Order("created_at DESC").
		Find(&topics).Error

	return topics, err
}

func (r *postRepository) GetTopic(ctx context.Context, topicID string) (*models.Topic, error) {
	var topic models.Topic
	err := r.db.WithContext(ctx).Where("id = ?", topicID).First(&topic).Error
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

// 收藏功能相关方法
func (r *postRepository) FavoritePost(ctx context.Context, userID string, postID string) error {
	// 检查是否已经收藏
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.PostFavorite{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // 已经收藏过了
	}

	// 创建收藏记录
	favorite := &models.PostFavorite{
		ID:        xid.New().String(),
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now(),
	}

	return r.db.WithContext(ctx).Create(favorite).Error
}

func (r *postRepository) UnfavoritePost(ctx context.Context, userID string, postID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&models.PostFavorite{}).Error
}

func (r *postRepository) GetFavoriteList(ctx context.Context, userID string, page, pageSize int32) ([]*models.Post, error) {
	var posts []*models.Post
	offset := (page - 1) * pageSize

	err := r.db.WithContext(ctx).
		Table("posts").
		Joins("JOIN post_favorites ON posts.id = post_favorites.post_id").
		Where("post_favorites.user_id = ?", userID).
		Offset(int(offset)).
		Limit(int(pageSize)).
		Order("post_favorites.created_at DESC").
		Find(&posts).Error

	return posts, err
}

// 评分功能相关方法
func (r *postRepository) RatePost(ctx context.Context, rating *models.PostRating) error {
	// 检查是否已经评分过
	var existingRating models.PostRating
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", rating.UserID, rating.PostID).
		First(&existingRating).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新评分
		rating.ID = xid.New().String()
		rating.CreatedAt = time.Now()
		return r.db.WithContext(ctx).Create(rating).Error
	} else if err != nil {
		return err
	} else {
		// 更新现有评分
		existingRating.Score = rating.Score
		existingRating.UpdatedAt = time.Now()
		return r.db.WithContext(ctx).Save(&existingRating).Error
	}
}

func (r *postRepository) GetScoreRanking(ctx context.Context, page, pageSize int32) ([]*models.Post, error) {
	var posts []*models.Post
	offset := (page - 1) * pageSize

	err := r.db.WithContext(ctx).
		Table("posts").
		Select("posts.*, AVG(post_ratings.score) as avg_score").
		Joins("LEFT JOIN post_ratings ON posts.id = post_ratings.post_id").
		Group("posts.id").
		Order("avg_score DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&posts).Error

	return posts, err
}

// 获取匿名头像信息
func (r *postRepository) GetAnonymousAvatar(ctx context.Context, avatarID string) (*models.AnonymousAvatar, error) {
	var avatar models.AnonymousAvatar
	err := r.db.WithContext(ctx).Where("id = ?", avatarID).First(&avatar).Error
	if err != nil {
		return nil, err
	}
	return &avatar, nil
}
