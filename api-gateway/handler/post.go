package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"hupu/kitex_gen/post"
	"hupu/kitex_gen/post/postservice"
	"hupu/shared/config"
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

	// 创建帖子服务客户端
	client, err := postservice.NewClient("post", client.WithHostPorts(config.GlobalConfig.Services.Post.Host+":"+config.GlobalConfig.Services.Post.Port))
	if err != nil {
		hlog.Errorf("Create post client error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "服务连接失败",
		})
		return
	}

	// 调用帖子服务
	resp, err := client.CreatePost(ctx, &req)
	if err != nil {
		hlog.Errorf("CreatePost error: %v", err)
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

	// 创建帖子服务客户端
	client, err := postservice.NewClient("post", client.WithHostPorts(config.GlobalConfig.Services.Post.Host+":"+config.GlobalConfig.Services.Post.Port))
	if err != nil {
		hlog.Errorf("Create post client error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "服务连接失败",
		})
		return
	}

	// 调用帖子服务
	req := &post.GetPostRequest{PostId: postID}
	resp, err := client.GetPost(ctx, req)
	if err != nil {
		hlog.Errorf("GetPost error: %v", err)
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

	// 创建帖子服务客户端
	client, err := postservice.NewClient("post", client.WithHostPorts(config.GlobalConfig.Services.Post.Host+":"+config.GlobalConfig.Services.Post.Port))
	if err != nil {
		hlog.Errorf("Create post client error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "服务连接失败",
		})
		return
	}

	// 构建请求
	req := &post.GetPostListRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	if userIDStr != "" {
		req.UserId = &userIDStr
	}

	// 调用帖子服务
	resp, err := client.GetPostList(ctx, req)
	if err != nil {
		hlog.Errorf("GetPostList error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取帖子列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
