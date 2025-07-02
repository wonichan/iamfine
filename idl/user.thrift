namespace go user

// 用户状态枚举
enum UserStatus {
    ACTIVE = 1,      // 正常
    BANNED = 2,      // 封禁
    DELETED = 3      // 已删除
}

// 关系状态枚举
enum RelationshipStatus {
    SINGLE = 1,         // 单身
    DATING = 2,         // 恋爱中
    MARRIED = 3,        // 已婚
    DIVORCED = 4,       // 离异
    COMPLICATED = 5     // 复杂
}

// 年龄段枚举
enum AgeGroup {
    POST_90 = 1,        // 90后
    POST_95 = 2,        // 95后
    POST_00 = 3,        // 00后
    POST_80 = 4,        // 80后
    OTHER = 99          // 其他
}

// 匿名马甲
struct AnonymousProfile {
    1: string id
    2: string user_id
    3: string anonymous_name
    4: string avatar_color
    5: bool is_active
    6: i64 created_at
}

struct User {
    1: string id
    2: string role
    3: string username
    4: string nickname
    5: string avatar
    6: string phone
    7: string email
    8: UserStatus status
    9: i64 created_at
    10: i64 updated_at
    11: optional string bio                    // 个人简介
    12: RelationshipStatus relationship_status // 感情状态
    13: AgeGroup age_group                     // 年龄段
    14: optional string location              // 所在地
    15: i32 post_count                        // 发帖数
    16: i32 comment_count                     // 评论数
    17: i32 like_count                        // 获赞数
    18: i32 collect_count                     // 收藏数
    19: double average_score                  // 平均得分
    20: i32 follower_count                    // 粉丝数
    21: i32 following_count                   // 关注数
    22: bool is_verified                      // 是否认证
    23: list<string> tags                     // 用户标签
    24: list<AnonymousProfile> anonymous_profiles // 匿名马甲列表
}

struct RegisterRequest {
    1: string username
    2: string password
    3: string nickname
    4: string phone
    5: string email
}

struct RegisterResponse {
    1: i32 code
    2: string message
    3: User user
}

struct LoginRequest {
    1: string username
    2: string password
}

struct LoginResponse {
    1: i32 code
    2: string message
    3: User user
}

struct GetUserRequest {
    1: string user_id
}

struct GetUserResponse {
    1: i32 code
    2: string message
    3: User user
}

struct CreateUserRequest {
    1: string username
    2: string nickname
    3: string phone
    4: string email
    5: string password
    6: optional RelationshipStatus relationship_status
    7: optional AgeGroup age_group
    8: optional string location
    9: optional string bio
}

struct UpdateUserRequest {
    1: string id
    2: optional string nickname
    3: optional string avatar
    4: optional string phone
    5: optional string email
    6: optional string bio
    7: optional RelationshipStatus relationship_status
    8: optional AgeGroup age_group
    9: optional string location
    10: optional list<string> tags
}

struct UpdateUserResponse {
    1: i32 code
    2: string message
    3: User user
}

// 关注相关请求响应
struct FollowUserRequest {
    1: string user_id
    2: string target_user_id
}

struct FollowUserResponse {
    1: i32 code
    2: string message
}

struct UnfollowUserRequest {
    1: string user_id
    2: string target_user_id
}

struct UnfollowUserResponse {
    1: i32 code
    2: string message
}

struct GetFollowersRequest {
    1: string user_id
    2: i32 page
    3: i32 page_size
}

struct GetFollowersResponse {
    1: i32 code
    2: string message
    3: list<User> users
    4: i32 total
}

struct GetFollowingRequest {
    1: string user_id
    2: i32 page
    3: i32 page_size
}

struct GetFollowingResponse {
    1: i32 code
    2: string message
    3: list<User> users
    4: i32 total
}

// 匿名马甲管理
struct CreateAnonymousProfileRequest {
    1: string user_id
    2: string anonymous_name
    3: string avatar_color
}

struct CreateAnonymousProfileResponse {
    1: i32 code
    2: string message
    3: AnonymousProfile profile
}

struct GetAnonymousProfilesRequest {
    1: string user_id
}

struct GetAnonymousProfilesResponse {
    1: i32 code
    2: string message
    3: list<AnonymousProfile> profiles
}

struct UpdateAnonymousProfileRequest {
    1: string profile_id
    2: optional string anonymous_name
    3: optional string avatar_color
    4: optional bool is_active
}

struct UpdateAnonymousProfileResponse {
    1: i32 code
    2: string message
}

// 用户统计
struct GetUserStatsRequest {
    1: string user_id
}

struct GetUserStatsResponse {
    1: i32 code
    2: string message
    3: i32 post_count
    4: i32 comment_count
    5: i32 like_count
    6: i32 collect_count
    7: double average_score
    8: i32 follower_count
    9: i32 following_count
}

service UserService {
    RegisterResponse Register(1: RegisterRequest req)
    LoginResponse Login(1: LoginRequest req)
    GetUserResponse GetUser(1: GetUserRequest req)
    UpdateUserResponse UpdateUser(1: UpdateUserRequest req)
    
    // 关注功能
    FollowUserResponse FollowUser(1: FollowUserRequest req)
    UnfollowUserResponse UnfollowUser(1: UnfollowUserRequest req)
    GetFollowersResponse GetFollowers(1: GetFollowersRequest req)
    GetFollowingResponse GetFollowing(1: GetFollowingRequest req)
    
    // 匿名马甲管理
    CreateAnonymousProfileResponse CreateAnonymousProfile(1: CreateAnonymousProfileRequest req)
    GetAnonymousProfilesResponse GetAnonymousProfiles(1: GetAnonymousProfilesRequest req)
    UpdateAnonymousProfileResponse UpdateAnonymousProfile(1: UpdateAnonymousProfileRequest req)
    
    // 用户统计
    GetUserStatsResponse GetUserStats(1: GetUserStatsRequest req)
}