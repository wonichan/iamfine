package handler

import (
	"context"
	"net/http"
	"strconv"

	"hupu/kitex_gen/follow"
	"hupu/kitex_gen/follow/followservice"
	"hupu/shared/config"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
)

// 关注用户
func Follow(ctx context.Context, c *app.RequestContext) {
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
	var req follow.FollowRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 设置关注者ID
	req.FollowerId = userID.(string)

	// 创建关注服务客户端
	client, err := followservice.NewClient("follow", client.WithHostPorts(config.GlobalConfig.Services.Follow.Host+":"+config.GlobalConfig.Services.Follow.Port))
	if err != nil {
		hlog.Errorf("Create follow client error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "服务连接失败",
		})
		return
	}

	// 调用关注服务
	resp, err := client.Follow(ctx, &req)
	if err != nil {
		hlog.Errorf("Follow error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "关注失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 取消关注
func Unfollow(ctx context.Context, c *app.RequestContext) {
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
	var req follow.UnfollowRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 设置关注者ID
	req.FollowerId = userID.(string)

	// 创建关注服务客户端
	client, err := followservice.NewClient("follow", client.WithHostPorts(config.GlobalConfig.Services.Follow.Host+":"+config.GlobalConfig.Services.Follow.Port))
	if err != nil {
		hlog.Errorf("Create follow client error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "服务连接失败",
		})
		return
	}

	// 调用关注服务
	resp, err := client.Unfollow(ctx, &req)
	if err != nil {
		hlog.Errorf("Unfollow error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "取消关注失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取关注列表
func GetFollowList(ctx context.Context, c *app.RequestContext) {
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

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 创建关注服务客户端
	client, err := followservice.NewClient("follow", client.WithHostPorts(config.GlobalConfig.Services.Follow.Host+":"+config.GlobalConfig.Services.Follow.Port))
	if err != nil {
		hlog.Errorf("Create follow client error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "服务连接失败",
		})
		return
	}

	// 调用关注服务
	req := &follow.GetFollowListRequest{
		UserId:   userID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}
	resp, err := client.GetFollowList(ctx, req)
	if err != nil {
		hlog.Errorf("GetFollowList error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取关注列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取粉丝列表
func GetFollowerList(ctx context.Context, c *app.RequestContext) {
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

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 创建关注服务客户端
	client, err := followservice.NewClient("follow", client.WithHostPorts(config.GlobalConfig.Services.Follow.Host+":"+config.GlobalConfig.Services.Follow.Port))
	if err != nil {
		hlog.Errorf("Create follow client error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "服务连接失败",
		})
		return
	}

	// 调用关注服务
	req := &follow.GetFollowerListRequest{
		UserId:   userID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}
	resp, err := client.GetFollowerList(ctx, req)
	if err != nil {
		hlog.Errorf("GetFollowerList error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取粉丝列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
