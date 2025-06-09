package handler

import (
	"context"
	"hupu/kitex_gen/comment"
	"hupu/shared/models"

	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type CommentHandler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewCommentHandler(db *gorm.DB, rdb *redis.Client) *CommentHandler {
	return &CommentHandler{
		db:  db,
		rdb: rdb,
	}
}

func (h *CommentHandler) CreateComment(ctx context.Context, req *comment.CreateCommentRequest) (*comment.CreateCommentResponse, error) {
	// 生成评论ID
	commentID := xid.New().String()

	// 创建评论
	newComment := models.Comment{
		ID:       commentID,
		PostID:   req.PostId,
		UserID:   req.UserId,
		Content:  req.Content,
		ParentID: req.ParentId,
	}

	err := h.db.Create(&newComment).Error
	if err != nil {
		return &comment.CreateCommentResponse{
			Code:    500,
			Message: "创建评论失败",
		}, err
	}

	return &comment.CreateCommentResponse{
		Code:    200,
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
	var comments []models.Comment
	query := h.db.Where("post_id = ?", req.PostId)

	// 如果指定了父评论ID，则只查询子评论
	if req.ParentId != nil {
		query = query.Where("parent_id = ?", *req.ParentId)
	} else {
		// 否则查询顶级评论
		query = query.Where("parent_id IS NULL OR parent_id = ''")
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	err := query.Offset(int(offset)).Limit(int(req.PageSize)).Order("created_at DESC").Find(&comments).Error
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
		Code:     200,
		Message:  "查询成功",
		Comments: commentList,
	}, nil
}
