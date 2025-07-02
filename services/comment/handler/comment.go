package handler

import (
	"context"
	"hupu/kitex_gen/comment"
	"hupu/services/comment/repository"
	"hupu/shared/models"

	"gorm.io/gorm"
)

type CommentHandler struct {
	db repository.CommentRepository
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{
		db: repository.NewCommentRepository(db),
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

	newComment, err := h.db.CreateComment(ctx, commentModel)
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
	comments, err := h.db.GetCommentList(ctx, req.PostId, req.ParentId, req.Page, req.PageSize)
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
		Code:    0,
		Message: "查询成功",
		Data: &comment.CommentListData{
			List:    commentList,
			Total:   int32(len(commentList)),
			HasMore: false,
		},
	}, nil
}

func (h *CommentHandler) DeleteComment(ctx context.Context, req *comment.DeleteCommentRequest) (*comment.DeleteCommentResponse, error) {
	err := h.db.DeleteComment(ctx, req.CommentId, req.UserId)
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
	commentModel, err := h.db.GetCommentDetail(ctx, req.CommentId)
	if err != nil {
		return &comment.GetCommentResponse{
			Code:    500,
			Message: "获取评论失败",
		}, err
	}

	response := h.convertToCommentResponse(commentModel)

	// 如果需要包含回复
	if req.IncludeReplies {
		replies, err := h.db.GetCommentList(ctx, commentModel.PostID, &commentModel.ID, 1, 100)
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
	comments, err := h.db.GetUserCommentList(ctx, req.UserId, req.Page, req.PageSize)
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

func (c *CommentHandler) RateComment(ctx context.Context, req *comment.RateCommentRequest) (r *comment.RateCommentResponse, err error) {
	panic("not implemented") // TODO: Implement
}

func (c *CommentHandler) GetUserCommentRating(ctx context.Context, req *comment.GetUserCommentRatingRequest) (r *comment.GetUserCommentRatingResponse, err error) {
	panic("not implemented") // TODO: Implement
}

func (c *CommentHandler) UpdateCommentRating(ctx context.Context, req *comment.UpdateCommentRatingRequest) (r *comment.UpdateCommentRatingResponse, err error) {
	panic("not implemented") // TODO: Implement
}

func (c *CommentHandler) DeleteCommentRating(ctx context.Context, req *comment.DeleteCommentRatingRequest) (r *comment.DeleteCommentRatingResponse, err error) {
	panic("not implemented") // TODO: Implement
}
