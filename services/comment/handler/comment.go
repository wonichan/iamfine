package handler

import (
	"context"
	"fmt"
	"hupu/kitex_gen/comment"
	"hupu/services/comment/repository"
	"hupu/shared/constants"
	"hupu/shared/middleware"
	"hupu/shared/models"
	"hupu/shared/utils"
	"strings"
)

type CommentHandler struct {
	db *repository.CommentRepository
}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{
		db: repository.NewCommentRepository(),
	}
}

func (h *CommentHandler) CreateComment(ctx context.Context, req *comment.CreateCommentRequest) (*comment.CreateCommentResponse, error) {
	// 参数验证
	if err := middleware.ValidateCommentCreation(req.PostId, req.UserId, req.Content); err != nil {
		return &comment.CreateCommentResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("参数验证失败: %s", err),
		}, nil
	}

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
		avatar, err := h.db.GetAnonymousAvatar(*req.AnonymousProfileId)
		if err == nil {
			commentModel.AnonymousName = &avatar.Name
			commentModel.AnonymousColor = &avatar.Color
		}
	}

	// 处理图片
	if req.Images != nil {
		commentModel.Images = models.StringArray(req.Images)
	}

	newComment, err := h.db.CreateComment(commentModel)
	if err != nil {
		return &comment.CreateCommentResponse{
			Code:    constants.CommentCreateFailCode,
			Message: fmt.Sprintf("创建评论失败: %s", err),
		}, nil
	}

	return &comment.CreateCommentResponse{
		Code:    constants.SuccessCode,
		Message: "创建成功",
		Comment: h.convertToCommentResponse(newComment),
	}, nil
}

func (h *CommentHandler) GetCommentList(ctx context.Context, req *comment.GetCommentListRequest) (*comment.GetCommentListResponse, error) {
	// 参数验证
	if req.PostId == "" {
		return &comment.GetCommentListResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to get comment list, post id is empty"),
		}, nil
	}

	comments, err := h.db.GetCommentList(req.PostId, req.ParentId, req.Page, req.PageSize)
	if err != nil {
		return &comment.GetCommentListResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("查询评论失败: %s", err),
		}, nil
	}

	// 转换为响应格式
	var commentList []*comment.Comment
	for _, c := range comments {
		commentList = append(commentList, h.convertToCommentResponse(c))
	}
	return &comment.GetCommentListResponse{
		Code:    constants.SuccessCode,
		Message: "查询成功",
		Data: &comment.CommentListData{
			List:    commentList,
			Total:   int32(len(commentList)),
			HasMore: false,
		},
	}, nil
}

func (h *CommentHandler) DeleteComment(ctx context.Context, req *comment.DeleteCommentRequest) (*comment.DeleteCommentResponse, error) {
	// 参数验证
	if req.CommentId == "" || req.UserId == "" {
		return &comment.DeleteCommentResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to delete comment, comment id or user id is empty"),
		}, nil
	}

	err := h.db.DeleteComment(req.CommentId, req.UserId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return &comment.DeleteCommentResponse{
				Code:    constants.CommentNotFoundCode,
				Message: fmt.Sprintf("failed to delete comment, err:%s", constants.GetErrorMessage(constants.CommentNotFoundCode)),
			}, nil
		}
		return &comment.DeleteCommentResponse{
			Code:    constants.CommentDeleteFailCode,
			Message: fmt.Sprintf("failed to delete comment, err:%s", err),
		}, nil
	}

	return &comment.DeleteCommentResponse{
		Code:    constants.SuccessCode,
		Message: "删除成功",
	}, nil
}

func (h *CommentHandler) GetComment(ctx context.Context, req *comment.GetCommentRequest) (*comment.GetCommentResponse, error) {
	// 参数验证
	if req.CommentId == "" {
		return &comment.GetCommentResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to get comment, comment id is empty"),
		}, nil
	}

	commentModel, err := h.db.GetCommentDetail(req.CommentId)
	if err != nil {
		if strings.Contains(fmt.Sprintf("%s", err), "not found") {
			return &comment.GetCommentResponse{
				Code:    constants.CommentNotFoundCode,
				Message: fmt.Sprintf("failed to get comment, err:%s", constants.GetErrorMessage(constants.CommentNotFoundCode)),
			}, nil
		}
		return &comment.GetCommentResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("failed to get comment, err:%s", err),
		}, nil
	}

	response := h.convertToCommentResponse(commentModel)

	// 如果需要包含回复
	if req.IncludeReplies {
		replies, err := h.db.GetCommentList(commentModel.PostID, &commentModel.ID, 1, 100)
		if err == nil {
			var replyList []*comment.Comment
			for _, reply := range replies {
				replyList = append(replyList, h.convertToCommentResponse(reply))
			}
			response.Replies = replyList
		}
	}

	return &comment.GetCommentResponse{
		Code:    constants.SuccessCode,
		Message: "获取成功",
		Comment: response,
	}, nil
}

func (h *CommentHandler) GetUserComments(ctx context.Context, req *comment.GetUserCommentsRequest) (*comment.GetUserCommentsResponse, error) {
	// 参数验证
	if req.UserId == "" {
		return &comment.GetUserCommentsResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to get user comments, user id is empty"),
		}, nil
	}

	comments, err := h.db.GetUserCommentList(req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &comment.GetUserCommentsResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("获取用户评论失败: %s", err),
		}, nil
	}

	// 转换为响应格式
	var commentList []*comment.Comment
	for _, c := range comments {
		commentList = append(commentList, h.convertToCommentResponse(c))
	}

	return &comment.GetUserCommentsResponse{
		Code:     constants.SuccessCode,
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
	// 参数验证
	if req.UserId == "" || req.CommentId == "" || req.Score < 1 || req.Score > 5 {
		return &comment.RateCommentResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to rate comment, user id or comment id is empty, score is %d", req.Score),
		}, nil
	}

	// TODO: Implement
	return &comment.RateCommentResponse{
		Code:    constants.SuccessCode,
		Message: "评分成功",
	}, nil
}

func (c *CommentHandler) GetUserCommentRating(ctx context.Context, req *comment.GetUserCommentRatingRequest) (r *comment.GetUserCommentRatingResponse, err error) {
	// 参数验证
	if req.UserId == "" || req.CommentId == "" {
		return &comment.GetUserCommentRatingResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to get user comment rating, user id or comment id is empty"),
		}, nil
	}
	// TODO: Implement
	return &comment.GetUserCommentRatingResponse{
		Code:    constants.SuccessCode,
		Message: "获取成功",
		Score:   utils.Float64Ptr(0),
		Comment: utils.StringPtr(""),
	}, nil
}

func (c *CommentHandler) UpdateCommentRating(ctx context.Context, req *comment.UpdateCommentRatingRequest) (r *comment.UpdateCommentRatingResponse, err error) {
	// 参数验证
	if req.UserId == "" || req.CommentId == "" || req.Score < 1 || req.Score > 5 {
		return &comment.UpdateCommentRatingResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to update comment rating, user id or comment id is empty, score is %d", req.Score),
		}, nil
	}

	// TODO: Implement
	return &comment.UpdateCommentRatingResponse{
		Code:         constants.SuccessCode,
		Message:      "更新成功",
		AverageScore: 0,
		TotalRatings: 0,
	}, nil
}

func (c *CommentHandler) DeleteCommentRating(ctx context.Context, req *comment.DeleteCommentRatingRequest) (r *comment.DeleteCommentRatingResponse, err error) {
	panic("not implemented") // TODO: Implement
}
