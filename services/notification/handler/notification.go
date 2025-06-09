package handler

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"gorm.io/gorm"
	"hupu/kitex_gen/notification"
	"hupu/shared/models"
)

type NotificationHandler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewNotificationHandler(db *gorm.DB, rdb *redis.Client) *NotificationHandler {
	return &NotificationHandler{
		db:  db,
		rdb: rdb,
	}
}

func (h *NotificationHandler) CreateNotification(ctx context.Context, req *notification.CreateNotificationRequest) (*notification.CreateNotificationResponse, error) {
	// 生成通知ID
	notificationID := xid.New().String()

	// 创建通知
	newNotification := models.Notification{
		ID:       notificationID,
		UserID:   req.UserId,
		Type:     req.Type,
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
			Type:      newNotification.Type,
			Title:     newNotification.Title,
			Content:   newNotification.Content,
			TargetId:  newNotification.TargetID,
			IsRead:    newNotification.IsRead,
			CreatedAt: newNotification.CreatedAt.Unix(),
			UpdatedAt: newNotification.UpdatedAt.Unix(),
		},
	}, nil
}

func (h *NotificationHandler) GetNotificationList(ctx context.Context, req *notification.GetNotificationListRequest) (*notification.GetNotificationListResponse, error) {
	var notifications []models.Notification
	offset := (req.Page - 1) * req.PageSize
	err := h.db.Where("user_id = ?", req.UserId).Offset(int(offset)).Limit(int(req.PageSize)).Order("created_at DESC").Find(&notifications).Error
	if err != nil {
		return &notification.GetNotificationListResponse{
			Code:    500,
			Message: "查询通知列表失败",
		}, err
	}

	// 转换为响应格式
	var notificationList []*notification.Notification
	for _, n := range notifications {
		notificationList = append(notificationList, &notification.Notification{
			Id:        n.ID,
			UserId:    n.UserID,
			Type:      n.Type,
			Title:     n.Title,
			Content:   n.Content,
			TargetId:  n.TargetID,
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt.Unix(),
			UpdatedAt: n.UpdatedAt.Unix(),
		})
	}

	return &notification.GetNotificationListResponse{
		Code:          200,
		Message:       "查询成功",
		Notifications: notificationList,
	}, nil
}
