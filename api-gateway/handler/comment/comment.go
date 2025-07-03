package comment

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/comment"
	"hupu/shared/constants"
	"hupu/shared/log"
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
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] GetCommentList request started", traceId)

	// 获取帖子ID参数
	postID, ok := common.ValidatePostIDParam(c, "postId")
	if !ok {
		common.RespondBadRequest(c, constants.MsgParamError)
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
	resp, err := handler.GetCommentClient().GetCommentList(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "GetCommentList", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetCommentList", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// CreateComment 创建评论
// POST /api/comments
func CreateComment(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] CreateComment request started", traceId)

	// 需要认证
	userID, ok := common.RequireAuth(c)
	if !ok {
		common.RespondUnauthorized(c)
		return
	}

	// 解析请求体
	var reqBody CreateCommentRequest
	if err := c.BindJSON(&reqBody); err != nil {
		common.RespondBadRequest(c, constants.MsgParamError)
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
	resp, err := handler.GetCommentClient().CreateComment(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "CreateComment", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "CreateComment", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// DeleteComment 删除评论
// DELETE /api/comments/{id}
func DeleteComment(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] DeleteComment request started", traceId)

	// 需要认证
	userID, ok := common.RequireAuth(c)
	if !ok {
		common.RespondUnauthorized(c)
		return
	}

	// 获取评论ID参数
	commentID, ok := common.ValidateCommentIDParam(c, "id")
	if !ok {
		common.RespondBadRequest(c, constants.MsgParamError)
		return
	}

	// 构建请求
	req := &comment.DeleteCommentRequest{
		UserId:    userID,
		CommentId: commentID,
	}

	// 调用评论服务
	resp, err := handler.GetCommentClient().DeleteComment(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "DeleteComment", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "DeleteComment", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}
