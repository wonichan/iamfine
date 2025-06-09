package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/kitex_gen/notification"
	"hupu/shared/utils"
)

// 获取通知列表
func GetNotificationList(ctx context.Context, c *app.RequestContext) {
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

	// 调用通知服务
	req := &notification.GetNotificationListRequest{
		UserId:   userID.(string),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}
	resp, err := notificationClient.GetNotificationList(ctx, req)
	if err != nil {
		utils.GetLogger().Errorf("GetNotificationList error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取通知列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
