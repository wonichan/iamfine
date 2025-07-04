package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"

	"hupu/shared/models"
	"hupu/shared/utils"
)

type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository 创建新的帖子数据库仓库
func NewPostRepository() *PostRepository {
	return &PostRepository{
		db: utils.GetDB(),
	}
}

// CreatePost 创建帖子
func (r *PostRepository) CreatePost(ctx context.Context, post *models.Post) error {
	// 生成帖子ID
	post.ID = xid.New().String()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	// 使用事务确保数据一致性
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建帖子
	if err := tx.Create(post).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 如果有话题ID，更新话题的帖子数量
	if post.TopicID != nil {
		if err := tx.Model(&models.Topic{}).Where("id = ?", post.TopicID).Update("post_count", gorm.Expr("post_count + ?", 1)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetPost 根据ID获取帖子
func (r *PostRepository) GetPost(ctx context.Context, id string) (*models.Post, error) {
	var post models.Post
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostList 获取帖子列表
func (r *PostRepository) GetPostList(ctx context.Context, page, pageSize int64) ([]*models.Post, error) {
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
func (r *PostRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	post.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(post).Error
}

// DeletePost 删除帖子
func (r *PostRepository) DeletePost(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Post{}).Error
}

// GetPostsByUserID 根据用户ID获取帖子列表
func (r *PostRepository) GetPostsByUserID(ctx context.Context, userID string, page, pageSize int64) ([]*models.Post, error) {
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
func (r *PostRepository) GetPostListByTopic(ctx context.Context, topicID string, page, pageSize int64) ([]*models.Post, error) {
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
func (r *PostRepository) GetPostListByCategory(ctx context.Context, category string, page, pageSize int64) ([]*models.Post, error) {
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

// GetPostListWithConditions 根据多个条件获取帖子列表
func (r *PostRepository) GetPostListWithConditions(ctx context.Context, conditions map[string]interface{}, page, pageSize int64, sortType string) ([]*models.Post, error) {
	var posts []*models.Post
	offset := (page - 1) * pageSize

	query := r.db.WithContext(ctx)

	// 添加查询条件
	for key, value := range conditions {
		if value != nil && value != "" {
			switch key {
			case "topic_id":
				query = query.Where("topic_id = ?", value)
			case "category":
				query = query.Where("category = ?", value)
			case "user_id":
				query = query.Where("user_id = ?", value)
			case "is_anonymous":
				query = query.Where("is_anonymous = ?", value)
			}
		}
	}

	// 设置排序
	switch sortType {
	case "hot":
		query = query.Order("like_count DESC, comment_count DESC, created_at DESC")
	case "score":
		query = query.Order("score DESC, created_at DESC")
	case "latest":
		fallthrough
	default:
		query = query.Order("created_at DESC")
	}

	err := query.
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&posts).Error

	return posts, err
}

// IncrementViewCount 增加浏览次数
func (r *PostRepository) IncrementViewCount(ctx context.Context, postID string) error {
	return r.db.WithContext(ctx).
		Model(&models.Post{}).
		Where("id = ?", postID).
		Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// 话题管理相关方法
func (r *PostRepository) CreateTopic(ctx context.Context, topic *models.Topic) error {
	topic.ID = xid.New().String()
	topic.CreatedAt = time.Now()
	topic.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(topic).Error
}

func (r *PostRepository) GetTopicList(ctx context.Context, page, pageSize int32) ([]*models.Topic, error) {
	var topics []*models.Topic
	offset := (page - 1) * pageSize

	err := r.db.WithContext(ctx).
		Offset(int(offset)).
		Limit(int(pageSize)).
		Order("created_at DESC").
		Find(&topics).Error

	return topics, err
}

func (r *PostRepository) GetTopic(ctx context.Context, topicID string) (*models.Topic, error) {
	var topic models.Topic
	err := r.db.WithContext(ctx).Where("id = ?", topicID).First(&topic).Error
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

// 收藏功能相关方法
func (r *PostRepository) FavoritePost(ctx context.Context, userID string, postID string) error {
	// 检查帖子是否存在
	var post models.Post
	err := r.db.WithContext(ctx).Where("id = ?", postID).First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("post not found")
		}
		return fmt.Errorf("failed to check post existence: %w", err)
	}

	// 检查是否已经收藏
	var count int64
	err = r.db.WithContext(ctx).
		Model(&models.PostFavorite{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check favorite existence: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("post already favorited")
	}

	// 使用事务创建收藏记录
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建收藏记录
	favorite := &models.PostFavorite{
		ID:        xid.New().String(),
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now(),
	}

	if err := tx.Create(favorite).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create favorite: %w", err)
	}

	return tx.Commit().Error
}

func (r *PostRepository) UnfavoritePost(ctx context.Context, userID string, postID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&models.PostFavorite{}).Error
}

func (r *PostRepository) GetFavoriteList(ctx context.Context, userID string, page, pageSize int32) ([]*models.Post, error) {
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
func (r *PostRepository) RatePost(ctx context.Context, rating *models.PostRating) error {
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

func (r *PostRepository) GetScoreRanking(ctx context.Context, page, pageSize int32) ([]*models.Post, error) {
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
func (r *PostRepository) GetAnonymousAvatar(ctx context.Context, avatarID string) (*models.AnonymousAvatar, error) {
	var avatar models.AnonymousAvatar
	err := r.db.WithContext(ctx).Where("id = ?", avatarID).First(&avatar).Error
	if err != nil {
		return nil, err
	}
	return &avatar, nil
}

// GetUserRating 获取用户对帖子的评分
func (r *PostRepository) GetUserRating(ctx context.Context, userID, postID string) (*models.PostRating, error) {
	var rating models.PostRating
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userID, postID).
		First(&rating).Error
	if err != nil {
		return nil, err
	}
	return &rating, nil
}

// UpdateRating 更新评分
func (r *PostRepository) UpdateRating(ctx context.Context, userID, postID string, score int32, comment *string) error {
	updateData := map[string]interface{}{
		"score":      score,
		"updated_at": time.Now(),
	}
	if comment != nil {
		updateData["comment"] = *comment
	}

	return r.db.WithContext(ctx).
		Model(&models.PostRating{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Updates(updateData).Error
}

// DeleteRating 删除评分
func (r *PostRepository) DeleteRating(ctx context.Context, userID, postID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&models.PostRating{}).Error
}

// GetPostRatingStats 获取帖子评分统计信息
func (r *PostRepository) GetPostRatingStats(ctx context.Context, postID string) (float64, int32, error) {
	var result struct {
		AverageScore float64
		TotalRatings int32
	}

	err := r.db.WithContext(ctx).
		Model(&models.PostRating{}).
		Select("AVG(score) as average_score, COUNT(*) as total_ratings").
		Where("post_id = ?", postID).
		Scan(&result).Error

	if err != nil {
		return 0, 0, err
	}

	return result.AverageScore, result.TotalRatings, nil
}
