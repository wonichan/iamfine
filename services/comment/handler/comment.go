package handler

import (
	"context"
	"hupu/kitex_gen/comment"
	"hupu/services/comment/repository"
	"hupu/shared/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type CommentHandler struct {
	db  repository.CommentRepository
	rdb repository.CommentRepository
}

func NewCommentHandler(db *gorm.DB, rdb *redis.Client) *CommentHandler {
	return &CommentHandler{
		db:  repository.NewCommentRepository(db),
		rdb: repository.NewCommentRedisRepo(rdb),
	}
}

func (h *CommentHandler) CreateComment(ctx context.Context, req *comment.CreateCommentRequest) (*comment.CreateCommentResponse, error) {
	// 创建评论
	commentModel := &models.Comment{
		PostID:   req.PostId,
		UserID:   req.UserId,
		Content:  req.Content,
		ParentID: req.ParentId,
	}

	newComment, err := h.rdb.CreateComment(ctx, commentModel)
	if err != nil {
		return &comment.CreateCommentResponse{
			Code:    500,
			Message: "创建评论失败",
		}, err
	}

	return &comment.CreateCommentResponse{
		Code:    0,
		Message: "创建成功",
		Comment: &comment.Comment{
			Id:        newComment.ID,
			PostId:    newComment.PostID,
			UserId:    newComment.UserID,
			Content:   newComment.Content,
			ParentId:  newComment.ParentID,
			LikeCount: newComment.LikeCount,
			CreatedAt: newComment.CreatedAt.Unix(),
			UpdatedAt: newComment.UpdatedAt.Unix(),
		},
	}, nil
}

func (h *CommentHandler) GetCommentList(ctx context.Context, req *comment.GetCommentListRequest) (*comment.GetCommentListResponse, error) {
	comments, err := h.rdb.GetCommentList(ctx, req.PostId, req.ParentId, req.Page, req.PageSize)
	if err != nil {
		return &comment.GetCommentListResponse{
			Code:    500,
			Message: "查询评论失败",
		}, err
	}

	// 转换为响应格式
	var commentList []*comment.Comment
	for _, c := range comments {
		commentList = append(commentList, &comment.Comment{
			Id:        c.ID,
			PostId:    c.PostID,
			UserId:    c.UserID,
			Content:   c.Content,
			ParentId:  c.ParentID,
			LikeCount: c.LikeCount,
			CreatedAt: c.CreatedAt.Unix(),
			UpdatedAt: c.UpdatedAt.Unix(),
		})
	}

	return &comment.GetCommentListResponse{
		Code:     0,
		Message:  "查询成功",
		Comments: commentList,
	}, nil
}
