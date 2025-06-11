package repository

import (
	"context"

	"hupu/shared/models"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetCommentList(ctx context.Context, postID string, parentID *string, page, pageSize int32) ([]*models.Comment, error)
	GetComment(ctx context.Context, commentID string) (*models.Comment, error)
	UpdateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	DeleteComment(ctx context.Context, commentID string) error
}