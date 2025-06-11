package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/kitex_gen/post"
	"hupu/shared/log"
)

// 创建帖子
func CreatePost(ctx context.Context, c *app.RequestContext) {
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
	var req post.CreatePostRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 设置用户ID
	req.UserId = userID.(string)

	// 调用帖子服务
	resp, err := postClient.CreatePost(ctx, &req)
	if err != nil {
		log.GetLogger().Errorf("CreatePost error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建帖子失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取帖子详情
func GetPost(ctx context.Context, c *app.RequestContext) {
	// 获取帖子ID参数
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "帖子ID不能为空",
		})
		return
	}

	// 调用帖子服务
	req := &post.GetPostRequest{PostId: postID}
	resp, err := postClient.GetPost(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetPost error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取帖子失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取帖子列表
func GetPostList(ctx context.Context, c *app.RequestContext) {
	// 解析查询参数
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	userIDStr := c.Query("user_id")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 构建请求
	req := &post.GetPostListRequest{
		Page:     page,
		PageSize: pageSize,
	}

	if userIDStr != "" {
		req.UserId = &userIDStr
	}

	// 调用帖子服务
	resp, err := postClient.GetPostList(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetPostList error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取帖子列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
