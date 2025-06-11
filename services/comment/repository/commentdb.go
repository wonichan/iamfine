package repository

import (
	"context"
	"fmt"

	"github.com/rs/xid"
	"gorm.io/gorm"

	"hupu/shared/models"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (cr *commentRepository) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	// 生成评论ID
	commentID := xid.New().String()

	// 创建评论
	newComment := &models.Comment{
		ID:       commentID,
		PostID:   comment.PostID,
		UserID:   comment.UserID,
		Content:  comment.Content,
		ParentID: comment.ParentID,
	}

	err := cr.db.Create(newComment).Error
	if err != nil {
		return nil, err
	}

	return newComment, nil
}

func (cr *commentRepository) GetCommentList(ctx context.Context, postID string, parentID *string, page, pageSize int32) ([]*models.Comment, error) {
	var comments []models.Comment
	query := cr.db.Where("post_id = ?", postID)

	// 如果指定了父评论ID，则只查询子评论
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		// 否则查询顶级评论
		query = query.Where("parent_id IS NULL OR parent_id = ''")
	}

	// 分页
	offset := (page - 1) * pageSize
	err := query.Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&comments).Error
	if err != nil {
		return nil, err
	}

	// 转换为指针切片
	var commentList []*models.Comment
	for i := range comments {
		commentList = append(commentList, &comments[i])
	}

	return commentList, nil
}

func (cr *commentRepository) GetComment(ctx context.Context, commentID string) (*models.Comment, error) {
	var comment models.Comment
	err := cr.db.Where("id = ?", commentID).First(&comment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("评论不存在")
		}
		return nil, err
	}
	return &comment, nil
}

func (cr *commentRepository) UpdateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	err := cr.db.Model(comment).Where("id = ?", comment.ID).Updates(comment).Error
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (cr *commentRepository) DeleteComment(ctx context.Context, commentID string) error {
	return cr.db.Where("id = ?", commentID).Delete(&models.Comment{}).Error
}