package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/kitex_gen/comment"
	"hupu/shared/log"
)

// 创建评论
func CreateComment(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 解析请求参数
	var req comment.CreateCommentRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 设置用户ID
	req.UserId = userID.(string)

	// 调用评论服务
	resp, err := commentClient.CreateComment(ctx, &req)
	if err != nil {
		log.GetLogger().Errorf("CreateComment error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建评论失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取评论列表
func GetCommentList(ctx context.Context, c *app.RequestContext) {
	// 获取帖子ID参数
	postID := c.Param("post_id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "帖子ID不能为空",
		})
		return
	}

	// 解析查询参数
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	sortType := c.Query("sort_type")
	parentID := c.Query("parent_id")
	includeRepliesStr := c.Query("include_replies")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	includeReplies := false
	if includeRepliesStr != "" {
		if parsed, err := strconv.ParseBool(includeRepliesStr); err == nil {
			includeReplies = parsed
		}
	}

	// 调用评论服务
	req := &comment.GetCommentListRequest{
		PostId:         postID,
		Page:           int32(page),
		PageSize:       int32(pageSize),
		IncludeReplies: includeReplies,
	}

	if sortType != "" {
		req.SortType = &sortType
	}

	if parentID != "" {
		req.ParentId = &parentID
	}

	resp, err := commentClient.GetCommentList(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetCommentList error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取评论列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取评论详情
func GetComment(ctx context.Context, c *app.RequestContext) {
	// 获取评论ID参数
	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "评论ID不能为空",
		})
		return
	}

	// 解析查询参数
	includeRepliesStr := c.Query("include_replies")
	includeReplies := false
	if includeRepliesStr != "" {
		if parsed, err := strconv.ParseBool(includeRepliesStr); err == nil {
			includeReplies = parsed
		}
	}

	// 构建请求
	req := &comment.GetCommentRequest{
		CommentId:      commentID,
		IncludeReplies: includeReplies,
	}

	// 调用评论服务
	resp, err := commentClient.GetComment(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetComment error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取评论详情失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取用户评论列表
func GetUserComments(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 解析查询参数
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	sortType := c.Query("sort_type")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 构建请求
	req := &comment.GetUserCommentsRequest{
		UserId:   userID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	if sortType != "" {
		req.SortType = &sortType
	}

	// 调用评论服务
	resp, err := commentClient.GetUserComments(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetUserComments error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取用户评论列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 点赞评论
func LikeComment(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取评论ID参数
	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "评论ID不能为空",
		})
		return
	}

	// 构建请求
	req := &comment.LikeCommentRequest{
		UserId:    userID.(string),
		CommentId: commentID,
	}

	// 调用评论服务
	resp, err := commentClient.LikeComment(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("LikeComment error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "点赞评论失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 取消点赞评论
func UnlikeComment(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取评论ID参数
	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "评论ID不能为空",
		})
		return
	}

	// 构建请求
	req := &comment.UnlikeCommentRequest{
		UserId:    userID.(string),
		CommentId: commentID,
	}

	// 调用评论服务
	resp, err := commentClient.UnlikeComment(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("UnlikeComment error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "取消点赞评论失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 删除评论
func DeleteComment(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取评论ID参数
	commentID := c.Param("id")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "评论ID不能为空",
		})
		return
	}

	// 构建请求
	req := &comment.DeleteCommentRequest{
		UserId:    userID.(string),
		CommentId: commentID,
	}

	// 调用评论服务
	resp, err := commentClient.DeleteComment(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("DeleteComment error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除评论失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
