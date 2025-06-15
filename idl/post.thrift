namespace go post

// 帖子分类枚举
enum PostCategory {
    DAILY_SHARE = 1,      // 日常分享
    LOVE_STORY = 2,       // 恋爱日常
    MARRIAGE_LIFE = 3,    // 婚姻围城
    FAMILY_RELATION = 4,  // 家庭关系
    EMOTIONAL_HELP = 5,   // 情感求助
    COMPLAINT = 6,        // 我要吐槽
    OTHER = 99           // 其他
}

// 话题结构
struct Topic {
    1: string id
    2: string name
    3: string description
    4: string icon
    5: string color
    6: i32 participant_count
    7: i64 created_at
    8: i64 updated_at
}

struct Post {
    1: string id
    2: string user_id
    3: string title
    4: string content
    5: list<string> images
    6: i32 like_count
    7: i32 comment_count
    8: double score
    9: i64 created_at
    10: i64 updated_at
    11: optional string topic_id        // 关联话题ID
    12: PostCategory category           // 帖子分类
    13: bool is_anonymous              // 是否匿名发布
    14: optional string anonymous_name          // 匿名用户名
    15: i32 view_count                 // 浏览次数
    16: i32 share_count                // 分享次数
    17: i32 collect_count              // 收藏次数
    18: bool is_hot                    // 是否热门
    19: bool is_top                    // 是否置顶
    20: optional string location       // 发布位置
    21: list<string> tags              // 标签列表
}

struct CreatePostRequest {
    1: string user_id
    2: string title
    3: string content
    4: list<string> images
    5: optional string topic_id        // 关联话题ID
    6: PostCategory category           // 帖子分类
    7: bool is_anonymous              // 是否匿名发布
    8: optional string location       // 发布位置
    9: list<string> tags              // 标签列表
}

struct CreatePostResponse {
    1: i32 code
    2: string message
    3: Post post
}

struct GetPostRequest {
    1: string post_id
}

struct GetPostResponse {
    1: i32 code
    2: string message
    3: Post post
}

struct GetPostListRequest {
    1: i64 page
    2: i64 page_size
    3: optional string user_id
    4: optional string topic_id       // 按话题筛选
    5: optional PostCategory category // 按分类筛选
    6: optional string sort_type      // 排序类型: latest, hot, score
    7: optional bool is_anonymous     // 是否只看匿名帖子
}

struct GetPostListResponse {
    1: i32 code
    2: string message
    3: list<Post> posts
    4: i32 total
}

// 话题相关请求响应
struct CreateTopicRequest {
    1: string name
    2: string description
    3: string icon
    4: string color
}

struct CreateTopicResponse {
    1: i32 code
    2: string message
    3: Topic topic
}

struct GetTopicListRequest {
    1: i32 page
    2: i32 page_size
    3: optional string sort_type      // hot, latest, participant
}

struct GetTopicListResponse {
    1: i32 code
    2: string message
    3: list<Topic> topics
    4: i32 total
}

// 收藏相关请求响应
struct CollectPostRequest {
    1: string user_id
    2: string post_id
}

struct CollectPostResponse {
    1: i32 code
    2: string message
}

struct UncollectPostRequest {
    1: string user_id
    2: string post_id
}

struct UncollectPostResponse {
    1: i32 code
    2: string message
}

struct GetCollectedPostsRequest {
    1: string user_id
    2: i32 page
    3: i32 page_size
}

struct GetCollectedPostsResponse {
    1: i32 code
    2: string message
    3: list<Post> posts
    4: i32 total
}

// 评分相关请求响应
struct RatePostRequest {
    1: string user_id
    2: string post_id
    3: double score              // 1-10分
    4: optional string comment   // 评分评论
}

struct RatePostResponse {
    1: i32 code
    2: string message
    3: double average_score      // 平均分
    4: i32 total_ratings        // 总评分数
}

// 获取评分排行榜
struct GetRatingRankRequest {
    1: i32 page
    2: i32 page_size
    3: string rank_type         // daily_high, daily_low, weekly_best, controversial
    4: optional string date     // 指定日期，格式: 2024-01-01
}

struct GetRatingRankResponse {
    1: i32 code
    2: string message
    3: list<Post> posts
    4: i32 total
}

service PostService {
    CreatePostResponse CreatePost(1: CreatePostRequest req)
    GetPostResponse GetPost(1: GetPostRequest req)
    GetPostListResponse GetPostList(1: GetPostListRequest req)
    
    // 话题管理
    CreateTopicResponse CreateTopic(1: CreateTopicRequest req)
    GetTopicListResponse GetTopicList(1: GetTopicListRequest req)
    
    // 收藏功能
    CollectPostResponse CollectPost(1: CollectPostRequest req)
    UncollectPostResponse UncollectPost(1: UncollectPostRequest req)
    GetCollectedPostsResponse GetCollectedPosts(1: GetCollectedPostsRequest req)
    
    // 评分功能
    RatePostResponse RatePost(1: RatePostRequest req)
    GetRatingRankResponse GetRatingRank(1: GetRatingRankRequest req)
}