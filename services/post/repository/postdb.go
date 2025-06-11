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
