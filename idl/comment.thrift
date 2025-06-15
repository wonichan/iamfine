namespace go comment

struct Comment {
    1: string id
    2: string post_id
    3: string user_id
    4: string content
    5: i64 created_at
    6: i64 updated_at
    7: optional string parent_id          // 父评论ID，用于回复
    8: i32 like_count                     // 点赞数
    9: i32 reply_count                    // 回复数
    10: bool is_anonymous                 // 是否匿名评论
    11: optional string anonymous_name    // 匿名用户名
    12: optional string anonymous_color   // 匿名用户颜色
    13: list<string> images               // 评论图片
    14: bool is_deleted                   // 是否已删除
    15: optional string location          // 评论位置
    16: list<Comment> replies             // 子评论列表（用于展示）
}

struct CreateCommentRequest {
    1: string post_id
    2: string user_id
    3: string content
    4: optional string parent_id          // 回复的评论ID
    5: bool is_anonymous                  // 是否匿名评论
    6: optional string anonymous_profile_id // 使用的匿名马甲ID
    7: list<string> images                // 评论图片
    8: optional string location           // 评论位置
}

struct CreateCommentResponse {
    1: i32 code
    2: string message
    3: Comment comment
}

struct GetCommentListRequest {
    1: string post_id
    2: i32 page
    3: i32 page_size
    4: optional string sort_type          // latest, hot, oldest
    5: optional string parent_id          // 获取特定评论的回复
    6: bool include_replies               // 是否包含回复
}

struct GetCommentListResponse {
    1: i32 code
    2: string message
    3: list<Comment> comments
    4: i32 total
}

// 点赞评论相关请求响应
struct LikeCommentRequest {
    1: string user_id
    2: string comment_id
}

struct LikeCommentResponse {
    1: i32 code
    2: string message
    3: i32 like_count
}

struct UnlikeCommentRequest {
    1: string user_id
    2: string comment_id
}

struct UnlikeCommentResponse {
    1: i32 code
    2: string message
    3: i32 like_count
}

// 删除评论
struct DeleteCommentRequest {
    1: string user_id
    2: string comment_id
}

struct DeleteCommentResponse {
    1: i32 code
    2: string message
}

// 获取评论详情
struct GetCommentRequest {
    1: string comment_id
    2: bool include_replies
}

struct GetCommentResponse {
    1: i32 code
    2: string message
    3: Comment comment
}

// 获取用户的评论列表
struct GetUserCommentsRequest {
    1: string user_id
    2: i32 page
    3: i32 page_size
    4: optional string sort_type
}

struct GetUserCommentsResponse {
    1: i32 code
    2: string message
    3: list<Comment> comments
    4: i32 total
}

service CommentService {
    CreateCommentResponse CreateComment(1: CreateCommentRequest req)
    GetCommentListResponse GetCommentList(1: GetCommentListRequest req)
    GetCommentResponse GetComment(1: GetCommentRequest req)
    GetUserCommentsResponse GetUserComments(1: GetUserCommentsRequest req)
    
    // 点赞功能
    LikeCommentResponse LikeComment(1: LikeCommentRequest req)
    UnlikeCommentResponse UnlikeComment(1: UnlikeCommentRequest req)
    
    // 删除功能
    DeleteCommentResponse DeleteComment(1: DeleteCommentRequest req)
}