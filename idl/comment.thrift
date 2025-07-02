namespace go comment

struct Author {
    1: string id
    2: string nickname
    3: string avatar
}

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
    17: bool is_liked                     // 当前用户是否已点赞
    18: optional Author author            // 作者信息
    19: double score                      // 评论评分
    20: i32 rating_count                  // 评分人数
    21: bool is_rated                     // 当前用户是否已评分
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
    4: optional string sort_type          // latest, hot, oldest, score_high, score_low
    5: optional string parent_id          // 获取特定评论的回复
    6: bool include_replies               // 是否包含回复
}

struct GetCommentListResponse {
    1: i32 code
    2: string message
    3: CommentListData data
}

struct CommentListData {
    1: list<Comment> list
    2: i32 total
    3: bool hasMore
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
    4: optional string sort_type          // latest, hot, oldest, score_high, score_low
}

struct GetUserCommentsResponse {
    1: i32 code
    2: string message
    3: list<Comment> comments
    4: i32 total
}

// 评论评分相关请求响应
struct RateCommentRequest {
    1: string user_id
    2: string comment_id
    3: double score              // 1-10分
    4: optional string comment   // 评分评论
}

struct RateCommentResponse {
    1: i32 code
    2: string message
    3: double average_score      // 平均分
    4: i32 total_ratings        // 总评分数
}

// 获取用户对评论的评分
struct GetUserCommentRatingRequest {
    1: string user_id
    2: string comment_id
}

struct GetUserCommentRatingResponse {
    1: i32 code
    2: string message
    3: optional double score     // 用户评分，如果未评分则为空
    4: optional string comment   // 评分评论
}

// 更新评论评分
struct UpdateCommentRatingRequest {
    1: string user_id
    2: string comment_id
    3: double score
    4: optional string comment
}

struct UpdateCommentRatingResponse {
    1: i32 code
    2: string message
    3: double average_score
    4: i32 total_ratings
}

// 删除评论评分
struct DeleteCommentRatingRequest {
    1: string user_id
    2: string comment_id
}

struct DeleteCommentRatingResponse {
    1: i32 code
    2: string message
    3: double average_score
    4: i32 total_ratings
}

service CommentService {
    CreateCommentResponse CreateComment(1: CreateCommentRequest req)
    GetCommentListResponse GetCommentList(1: GetCommentListRequest req)
    GetCommentResponse GetComment(1: GetCommentRequest req)
    GetUserCommentsResponse GetUserComments(1: GetUserCommentsRequest req)
    
    // 删除功能
    DeleteCommentResponse DeleteComment(1: DeleteCommentRequest req)
    
    // 评论评分功能
    RateCommentResponse RateComment(1: RateCommentRequest req)
    GetUserCommentRatingResponse GetUserCommentRating(1: GetUserCommentRatingRequest req)
    UpdateCommentRatingResponse UpdateCommentRating(1: UpdateCommentRatingRequest req)
    DeleteCommentRatingResponse DeleteCommentRating(1: DeleteCommentRatingRequest req)
}