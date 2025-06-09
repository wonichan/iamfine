namespace go notification

struct Notification {
    1: string id
    2: string user_id
    3: string title
    4: string content
    5: i32 type // 1: like, 2: comment, 3: follow
    6: string target_id
    7: bool is_read
    8: i64 created_at
}

struct CreateNotificationRequest {
    1: string user_id
    2: string title
    3: string content
    4: i32 type
    5: string target_id
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

service NotificationService {
    CreateNotificationResponse CreateNotification(1: CreateNotificationRequest req)
    GetNotificationListResponse GetNotificationList(1: GetNotificationListRequest req)
}