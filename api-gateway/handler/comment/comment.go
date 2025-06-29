package comment

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/comment"
	"hupu/shared/constants"
)

// CreateCommentRequest 创建评论请求结构
type CreateCommentRequest struct {
	PostID   string  `json:"postId" binding:"required"`
	ParentID *string `json:"parentId,omitempty"`
	Content  string  `json:"content" binding:"required"`
}

// GetCommentList 获取评论列表
// GET /api/comments
func GetCommentList(ctx context.Context, c *app.RequestContext) {
	// 获取帖子ID参数
	postID, ok := common.ValidatePostIDParam(c, "postId")
	if !ok {
		return
	}

	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 解析可选参数
	sortType := common.ParseOptionalStringParam(c, "sort_type")
	parentID := common.ParseOptionalStringParam(c, "parent_id")
	includeReplies := common.ParseOptionalBoolParam(c, "include_replies")

	// 构建请求
	req := &comment.GetCommentListRequest{
		PostId:         postID,
		Page:           page,
		PageSize:       pageSize,
		IncludeReplies: includeReplies != nil && *includeReplies,
	}

	if sortType != nil {
		req.SortType = sortType
	}

	if parentID != nil {
		req.ParentId = parentID
	}

	// 调用评论服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetCommentClient().GetCommentList(ctx, req)
	}), "GetCommentList", constants.MsgGetCommentListFailed)
}

// CreateComment 创建评论
// POST /api/comments
func CreateComment(ctx context.Context, c *app.RequestContext) {
	// 需要认证
	userID, ok := common.RequireAuth(c)
	if !ok {
		return
	}

	// 解析请求体
	var reqBody CreateCommentRequest
	if err := c.BindJSON(&reqBody); err != nil {
		common.ErrorResponseFunc(c, constants.HTTPStatusBadRequest, common.CodeError, constants.MsgParamError)
		return
	}

	// 构建请求
	req := &comment.CreateCommentRequest{
		PostId:  reqBody.PostID,
		UserId:  userID,
		Content: reqBody.Content,
	}

	if reqBody.ParentID != nil {
		req.ParentId = reqBody.ParentID
	}

	// 调用评论服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetCommentClient().CreateComment(ctx, req)
	}), "CreateComment", constants.MsgCreateCommentFailed)
}

// DeleteComment 删除评论
// DELETE /api/comments/{id}
func DeleteComment(ctx context.Context, c *app.RequestContext) {
	// 需要认证
	userID, ok := common.RequireAuth(c)
	if !ok {
		return
	}

	// 获取评论ID参数
	commentID, ok := common.ValidateCommentIDParam(c, "id")
	if !ok {
		return
	}

	// 构建请求
	req := &comment.DeleteCommentRequest{
		UserId:    userID,
		CommentId: commentID,
	}

	// 调用评论服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetCommentClient().DeleteComment(ctx, req)
	}), "DeleteComment", constants.MsgDeleteCommentFailed)
}

// LikeComment 点赞评论
// POST /api/comments/{id}/like
func LikeComment(ctx context.Context, c *app.RequestContext) {
	// 需要认证
	userID, ok := common.RequireAuth(c)
	if !ok {
		return
	}

	// 获取评论ID参数
	commentID, ok := common.ValidateCommentIDParam(c, "id")
	if !ok {
		return
	}

	// 构建请求
	req := &comment.LikeCommentRequest{
		UserId:    userID,
		CommentId: commentID,
	}

	// 调用评论服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetCommentClient().LikeComment(ctx, req)
	}), "LikeComment", constants.MsgLikeCommentFailed)
}

// UnlikeComment 取消点赞评论
// DELETE /api/comments/{id}/like
func UnlikeComment(ctx context.Context, c *app.RequestContext) {
	// 需要认证
	userID, ok := common.RequireAuth(c)
	if !ok {
		return
	}

	// 获取评论ID参数
	commentID, ok := common.ValidateCommentIDParam(c, "id")
	if !ok {
		return
	}

	// 构建请求
	req := &comment.UnlikeCommentRequest{
		UserId:    userID,
		CommentId: commentID,
	}

	// 调用评论服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetCommentClient().UnlikeComment(ctx, req)
	}), "UnlikeComment", constants.MsgUnlikeCommentFailed)
}
