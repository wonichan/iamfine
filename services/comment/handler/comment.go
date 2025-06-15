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
		PostID:      req.PostId,
		UserID:      req.UserId,
		Content:     req.Content,
		ParentID:    req.ParentId,
		IsAnonymous: req.IsAnonymous,
		Location:    req.Location,
	}

	// 处理匿名评论
	if req.IsAnonymous && req.AnonymousProfileId != nil {
		// 根据匿名马甲ID获取匿名信息
		avatar, err := h.db.GetAnonymousAvatar(ctx, *req.AnonymousProfileId)
		if err == nil {
			commentModel.AnonymousName = &avatar.Name
			commentModel.AnonymousColor = &avatar.Color
		}
	}

	// 处理图片
	if req.Images != nil {
		commentModel.Images = models.StringArray(req.Images)
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
		Comment: h.convertToCommentResponse(newComment),
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
		commentList = append(commentList, h.convertToCommentResponse(c))
	}

	return &comment.GetCommentListResponse{
		Code:     0,
		Message:  "查询成功",
		Comments: commentList,
	}, nil
}

func (h *CommentHandler) LikeComment(ctx context.Context, req *comment.LikeCommentRequest) (*comment.LikeCommentResponse, error) {
	err := h.rdb.LikeComment(ctx, req.CommentId, req.UserId)
	if err != nil {
		return &comment.LikeCommentResponse{
			Code:    500,
			Message: "点赞失败",
		}, err
	}

	return &comment.LikeCommentResponse{
		Code:    0,
		Message: "点赞成功",
	}, nil
}

func (h *CommentHandler) UnlikeComment(ctx context.Context, req *comment.UnlikeCommentRequest) (*comment.UnlikeCommentResponse, error) {
	err := h.rdb.UnlikeComment(ctx, req.CommentId, req.UserId)
	if err != nil {
		return &comment.UnlikeCommentResponse{
			Code:    500,
			Message: "取消点赞失败",
		}, err
	}

	return &comment.UnlikeCommentResponse{
		Code:    0,
		Message: "取消点赞成功",
	}, nil
}

func (h *CommentHandler) DeleteComment(ctx context.Context, req *comment.DeleteCommentRequest) (*comment.DeleteCommentResponse, error) {
	err := h.rdb.DeleteComment(ctx, req.CommentId, req.UserId)
	if err != nil {
		return &comment.DeleteCommentResponse{
			Code:    500,
			Message: "删除失败",
		}, err
	}

	return &comment.DeleteCommentResponse{
		Code:    0,
		Message: "删除成功",
	}, nil
}

func (h *CommentHandler) GetComment(ctx context.Context, req *comment.GetCommentRequest) (*comment.GetCommentResponse, error) {
	commentModel, err := h.rdb.GetCommentDetail(ctx, req.CommentId)
	if err != nil {
		return &comment.GetCommentResponse{
			Code:    500,
			Message: "获取评论失败",
		}, err
	}

	response := h.convertToCommentResponse(commentModel)

	// 如果需要包含回复
	if req.IncludeReplies {
		replies, err := h.rdb.GetCommentList(ctx, commentModel.PostID, &commentModel.ID, 1, 100)
		if err == nil {
			var replyList []*comment.Comment
			for _, reply := range replies {
				replyList = append(replyList, h.convertToCommentResponse(reply))
			}
			response.Replies = replyList
		}
	}

	return &comment.GetCommentResponse{
		Code:    0,
		Message: "获取成功",
		Comment: response,
	}, nil
}

func (h *CommentHandler) GetUserComments(ctx context.Context, req *comment.GetUserCommentsRequest) (*comment.GetUserCommentsResponse, error) {
	comments, err := h.rdb.GetUserCommentList(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &comment.GetUserCommentsResponse{
			Code:    500,
			Message: "获取用户评论失败",
		}, err
	}

	// 转换为响应格式
	var commentList []*comment.Comment
	for _, c := range comments {
		commentList = append(commentList, h.convertToCommentResponse(c))
	}

	return &comment.GetUserCommentsResponse{
		Code:     0,
		Message:  "获取成功",
		Comments: commentList,
		Total:    int32(len(commentList)),
	}, nil
}

// 辅助方法：转换模型为响应格式
func (h *CommentHandler) convertToCommentResponse(c *models.Comment) *comment.Comment {
	response := &comment.Comment{
		Id:          c.ID,
		PostId:      c.PostID,
		UserId:      c.UserID,
		Content:     c.Content,
		ParentId:    c.ParentID,
		LikeCount:   c.LikeCount,
		ReplyCount:  c.ReplyCount,
		IsAnonymous: c.IsAnonymous,
		IsDeleted:   c.IsDeleted,
		CreatedAt:   c.CreatedAt.Unix(),
		UpdatedAt:   c.UpdatedAt.Unix(),
	}

	if c.AnonymousName != nil {
		response.AnonymousName = c.AnonymousName
	}
	if c.AnonymousColor != nil {
		response.AnonymousColor = c.AnonymousColor
	}
	if c.Location != nil {
		response.Location = c.Location
	}
	if len(c.Images) > 0 {
		images := make([]string, len(c.Images))
		copy(images, c.Images)
		response.Images = images
	}

	return response
}
