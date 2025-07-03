package handler

import (
	"context"
	"hupu/kitex_gen/notification"
	"hupu/shared/models"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type NotificationHandler struct {
	db *gorm.DB
}

func NewNotificationHandler(db *gorm.DB) *NotificationHandler {
	return &NotificationHandler{
		db: db,
	}
}

func (h *NotificationHandler) CreateNotification(ctx context.Context, req *notification.CreateNotificationRequest) (*notification.CreateNotificationResponse, error) {
	// 生成通知ID
	notificationID := xid.New().String()

	// 创建通知
	newNotification := models.Notification{
		ID:       notificationID,
		UserID:   req.UserId,
		Type:     int32(req.Type),
		Title:    req.Title,
		Content:  req.Content,
		TargetID: req.TargetId,
		IsRead:   false,
	}

	err := h.db.Create(&newNotification).Error
	if err != nil {
		return &notification.CreateNotificationResponse{
			Code:    500,
			Message: "创建通知失败",
		}, err
	}

	return &notification.CreateNotificationResponse{
		Code:    200,
		Message: "创建成功",
		Notification: &notification.Notification{
			Id:        newNotification.ID,
			UserId:    newNotification.UserID,
			Type:      notification.NotificationType(newNotification.Type),
			Title:     newNotification.Title,
			Content:   newNotification.Content,
			TargetId:  newNotification.TargetID,
			IsRead:    newNotification.IsRead,
			CreatedAt: newNotification.CreatedAt.Unix(),
		},
	}, nil
}

func (h *NotificationHandler) GetNotificationList(ctx context.Context, req *notification.GetNotificationListRequest) (*notification.GetNotificationListResponse, error) {
	var notifications []models.Notification
	offset := (req.Page - 1) * req.PageSize

	query := h.db.Where("user_id = ?", req.UserId)

	err := query.Offset(int(offset)).Limit(int(req.PageSize)).Order("created_at DESC").Find(&notifications).Error
	if err != nil {
		return &notification.GetNotificationListResponse{
			Code:    500,
			Message: "查询通知列表失败",
		}, err
	}

	// 获取总数
	var total int64
	h.db.Model(&models.Notification{}).Where("user_id = ?", req.UserId).Count(&total)

	// 转换为响应格式
	var notificationList []*notification.Notification
	for _, n := range notifications {
		notificationList = append(notificationList, h.convertToNotificationResponse(&n))
	}

	return &notification.GetNotificationListResponse{
		Code:          200,
		Message:       "查询成功",
		Notifications: notificationList,
		Total:         int32(total),
	}, nil
}

// MarkNotificationRead 标记通知为已读
func (h *NotificationHandler) MarkNotificationRead(ctx context.Context, req *notification.MarkNotificationReadRequest) (*notification.MarkNotificationReadResponse, error) {
	err := h.db.Model(&models.Notification{}).Where("id = ? AND user_id = ?", req.NotificationId, req.UserId).Update("is_read", true).Error
	if err != nil {
		return &notification.MarkNotificationReadResponse{
			Code:    500,
			Message: "标记已读失败",
		}, err
	}

	return &notification.MarkNotificationReadResponse{
		Code:    200,
		Message: "标记成功",
	}, nil
}

// MarkAllNotificationsRead 标记所有通知为已读
func (h *NotificationHandler) MarkAllNotificationsRead(ctx context.Context, req *notification.MarkAllNotificationsReadRequest) (*notification.MarkAllNotificationsReadResponse, error) {
	query := h.db.Model(&models.Notification{}).Where("user_id = ?", req.UserId)

	// 如果指定了类型，只标记特定类型的通知
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
	}

	// 执行更新并获取影响的行数
	result := query.Update("is_read", true)
	if result.Error != nil {
		return &notification.MarkAllNotificationsReadResponse{
			Code:    500,
			Message: "标记失败",
		}, result.Error
	}

	return &notification.MarkAllNotificationsReadResponse{
		Code:          200,
		Message:       "标记成功",
		AffectedCount: int32(result.RowsAffected),
	}, nil
}

// GetUnreadCount 获取未读通知数量
func (h *NotificationHandler) GetUnreadCount(ctx context.Context, req *notification.GetUnreadCountRequest) (*notification.GetUnreadCountResponse, error) {
	var count int64
	query := h.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", req.UserId, false)

	// 如果指定了类型，只统计特定类型的通知
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
	}

	err := query.Count(&count).Error
	if err != nil {
		return &notification.GetUnreadCountResponse{
			Code:    500,
			Message: "查询失败",
		}, err
	}

	return &notification.GetUnreadCountResponse{
		Code:        200,
		Message:     "查询成功",
		UnreadCount: int32(count),
	}, nil
}

// DeleteNotification 删除通知
func (h *NotificationHandler) DeleteNotification(ctx context.Context, req *notification.DeleteNotificationRequest) (*notification.DeleteNotificationResponse, error) {
	err := h.db.Where("id = ? AND user_id = ?", req.NotificationId, req.UserId).Delete(&models.Notification{}).Error
	if err != nil {
		return &notification.DeleteNotificationResponse{
			Code:    500,
			Message: "删除失败",
		}, err
	}

	return &notification.DeleteNotificationResponse{
		Code:    200,
		Message: "删除成功",
	}, nil
}

// 辅助方法：转换模型为响应格式
func (h *NotificationHandler) convertToNotificationResponse(n *models.Notification) *notification.Notification {
	return &notification.Notification{
		Id:        n.ID,
		UserId:    n.UserID,
		Type:      notification.NotificationType(n.Type),
		Title:     n.Title,
		Content:   n.Content,
		TargetId:  n.TargetID,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt.Unix(),
	}
}
