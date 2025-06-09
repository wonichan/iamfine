package handler

import (
	"context"
	"net/http"
	"strconv" // 添加这行

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/kitex_gen/like"
	"hupu/shared/log"
)

// 点赞
func Like(ctx context.Context, c *app.RequestContext) {
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
	var req like.LikeRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 设置用户ID
	req.UserId = userID.(string)

	// 调用点赞服务
	resp, err := likeClient.Like(ctx, &req)
	if err != nil {
		log.GetLogger().Errorf("Like error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "点赞失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 取消点赞
func Unlike(ctx context.Context, c *app.RequestContext) {
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
	var req like.UnlikeRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 设置用户ID
	req.UserId = userID.(string)

	// 调用点赞服务
	resp, err := likeClient.Unlike(ctx, &req)
	if err != nil {
		log.GetLogger().Errorf("Unlike error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "取消点赞失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取点赞列表
func GetLikeList(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
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

	// 调用点赞服务
	req := &like.GetLikeListRequest{
		UserId:   userID.(string),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}
	resp, err := likeClient.GetLikeList(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetLikeList error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取点赞列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
