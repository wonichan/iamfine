package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/kitex_gen/follow"
	"hupu/shared/log"
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

	// 调用关注服务
	resp, err := followClient.Follow(ctx, &req)
	if err != nil {
		log.GetLogger().Errorf("Follow error: %v", err)
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

	// 调用关注服务
	resp, err := followClient.Unfollow(ctx, &req)
	if err != nil {
		log.GetLogger().Errorf("Unfollow error: %v", err)
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

	// 调用关注服务
	req := &follow.GetFollowListRequest{
		UserId:   userID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}
	resp, err := followClient.GetFollowList(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetFollowList error: %v", err)
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

	// 调用关注服务
	req := &follow.GetFollowerListRequest{
		UserId:   userID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}
	resp, err := followClient.GetFollowerList(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetFollowerList error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取粉丝列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 检查关注状态
func CheckFollowStatus(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取目标用户ID参数
	targetUserID := c.Query("target_user_id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "目标用户ID不能为空",
		})
		return
	}

	// 构建请求
	req := &follow.CheckFollowStatusRequest{
		FollowerId:  userID.(string),
		FollowingId: targetUserID,
	}

	// 调用关注服务
	resp, err := followClient.CheckFollowStatus(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("CheckFollowStatus error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "检查关注状态失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 检查是否关注
func IsFollowing(ctx context.Context, c *app.RequestContext) {
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

	// 调用关注服务
	resp, err := followClient.IsFollowing(ctx, &req)
	if err != nil {
		log.GetLogger().Errorf("IsFollowing error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "检查关注状态失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取关注数量
func GetFollowCount(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 构建请求
	req := &follow.GetFollowCountRequest{
		UserId: userID,
	}

	// 调用关注服务
	resp, err := followClient.GetFollowCount(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetFollowCount error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取关注数量失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取粉丝数量
func GetFollowerCount(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 构建请求
	req := &follow.GetFollowerCountRequest{
		UserId: userID,
	}

	// 调用关注服务
	resp, err := followClient.GetFollowerCount(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetFollowerCount error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取粉丝数量失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取共同关注
func GetMutualFollows(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取目标用户ID参数
	targetUserID := c.Query("target_user_id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "目标用户ID不能为空",
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

	// 构建请求
	req := &follow.GetMutualFollowsRequest{
		UserId:       userID.(string),
		TargetUserId: targetUserID,
		Page:         int32(page),
		PageSize:     int32(pageSize),
	}

	// 调用关注服务
	resp, err := followClient.GetMutualFollows(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetMutualFollows error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取共同关注失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
