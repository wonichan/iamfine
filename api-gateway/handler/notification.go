package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/kitex_gen/notification"
	"hupu/shared/log"
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
		log.GetLogger().Errorf("GetNotificationList error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取通知列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 标记通知为已读
func MarkNotificationRead(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取通知ID参数
	notificationID := c.Param("notification_id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "通知ID不能为空",
		})
		return
	}

	// 调用通知服务
	req := &notification.MarkNotificationReadRequest{
		UserId:         userID.(string),
		NotificationId: notificationID,
	}
	resp, err := notificationClient.MarkNotificationRead(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("MarkNotificationRead error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "标记已读失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 标记所有通知为已读
func MarkAllNotificationsRead(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 解析请求参数（可选的通知类型）
	var req notification.MarkAllNotificationsReadRequest
	if err := c.BindAndValidate(&req); err != nil {
		// 如果绑定失败，使用基本请求
		req = notification.MarkAllNotificationsReadRequest{
			UserId: userID.(string),
		}
	} else {
		req.UserId = userID.(string)
	}

	// 调用通知服务
	resp, err := notificationClient.MarkAllNotificationsRead(ctx, &req)
	if err != nil {
		log.GetLogger().Errorf("MarkAllNotificationsRead error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "标记失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取未读通知数量
func GetUnreadCount(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 解析查询参数（可选的通知类型）
	typeStr := c.Query("type")
	req := &notification.GetUnreadCountRequest{
		UserId: userID.(string),
	}

	// 如果指定了类型参数
	if typeStr != "" {
		typeInt, err := strconv.ParseInt(typeStr, 10, 32)
		if err == nil {
			notificationType := notification.NotificationType(typeInt)
			req.Type = &notificationType
		}
	}

	// 调用通知服务
	resp, err := notificationClient.GetUnreadCount(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetUnreadCount error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取未读数量失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 删除通知
func DeleteNotification(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取通知ID参数
	notificationID := c.Param("notification_id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "通知ID不能为空",
		})
		return
	}

	// 调用通知服务
	req := &notification.DeleteNotificationRequest{
		UserId:         userID.(string),
		NotificationId: notificationID,
	}
	resp, err := notificationClient.DeleteNotification(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("DeleteNotification error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除通知失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
