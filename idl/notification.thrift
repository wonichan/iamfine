namespace go notification

// 通知类型枚举
enum NotificationType {
    LIKE = 1,           // 点赞通知
    COMMENT = 2,        // 评论通知
    FOLLOW = 3,         // 关注通知
    REPLY = 4,          // 回复通知
    COLLECT = 5,        // 收藏通知
    RATE = 6,           // 评分通知
    SYSTEM = 7,         // 系统通知
    TOPIC_UPDATE = 8    // 话题更新通知
}

struct Notification {
    1: string id
    2: string user_id
    3: string title
    4: string content
    5: NotificationType type
    6: string target_id
    7: bool is_read
    8: i64 created_at
    9: optional string sender_id        // 发送者ID
    10: optional string sender_name     // 发送者名称
    11: optional string sender_avatar   // 发送者头像
    12: optional string extra_data      // 额外数据（JSON格式）
}

struct CreateNotificationRequest {
    1: string user_id
    2: string title
    3: string content
    4: NotificationType type
    5: string target_id
    6: optional string sender_id
    7: optional string extra_data
}

struct CreateNotificationResponse {
    1: i32 code
    2: string message
    3: Notification notification
}

struct GetNotificationListRequest {
    1: string user_id
    2: i32 page
    3: i32 page_size
}

struct GetNotificationListResponse {
    1: i32 code
    2: string message
    3: list<Notification> notifications
    4: i32 total
}

// 标记通知已读
struct MarkNotificationReadRequest {
    1: string user_id
    2: string notification_id
}

struct MarkNotificationReadResponse {
    1: i32 code
    2: string message
}

// 批量标记已读
struct MarkAllNotificationsReadRequest {
    1: string user_id
    2: optional NotificationType type  // 可选：只标记特定类型的通知
}

struct MarkAllNotificationsReadResponse {
    1: i32 code
    2: string message
    3: i32 affected_count
}

// 获取未读通知数量
struct GetUnreadCountRequest {
    1: string user_id
    2: optional NotificationType type
}

struct GetUnreadCountResponse {
    1: i32 code
    2: string message
    3: i32 unread_count
}

// 删除通知
struct DeleteNotificationRequest {
    1: string user_id
    2: string notification_id
}

struct DeleteNotificationResponse {
    1: i32 code
    2: string message
}

service NotificationService {
    CreateNotificationResponse CreateNotification(1: CreateNotificationRequest req)
    GetNotificationListResponse GetNotificationList(1: GetNotificationListRequest req)
    MarkNotificationReadResponse MarkNotificationRead(1: MarkNotificationReadRequest req)
    MarkAllNotificationsReadResponse MarkAllNotificationsRead(1: MarkAllNotificationsReadRequest req)
    GetUnreadCountResponse GetUnreadCount(1: GetUnreadCountRequest req)
    DeleteNotificationResponse DeleteNotification(1: DeleteNotificationRequest req)
}