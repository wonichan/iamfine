namespace go follow

// 用户基本信息结构体
struct UserInfo {
    1: string id
    2: string nickname
    3: string avatar
}

// 关注关系结构体
struct Follow {
    1: string id
    2: string follower_id      // 关注者ID
    3: string following_id     // 被关注者ID
    4: i64 created_at
    5: bool is_mutual          // 是否互相关注
    6: optional UserInfo user_info  // 用户信息
}

struct FollowRequest {
    1: string follower_id
    2: string following_id
}

struct FollowResponse {
    1: i32 code
    2: string message
    3: bool is_following
}

struct UnfollowRequest {
    1: string follower_id
    2: string following_id
}

struct UnfollowResponse {
    1: i32 code
    2: string message
}

// 检查关注状态
struct CheckFollowStatusRequest {
    1: string follower_id
    2: string following_id
}

struct CheckFollowStatusResponse {
    1: i32 code
    2: string message
    3: bool is_following
    4: bool is_mutual
}

struct GetFollowListRequest {
    1: string user_id
    2: i32 type // 1: following, 2: followers
    3: i32 page
    4: i32 page_size
}

struct GetFollowListResponse {
    1: i32 code
    2: string message
    3: list<Follow> follows
    4: i32 total
    5: bool has_more
}

// 添加 GetFollowerListRequest 和 GetFollowerListResponse
struct GetFollowerListRequest {
    1: string user_id
    2: i32 page
    3: i32 page_size
}

struct GetFollowerListResponse {
    1: i32 code
    2: string message
    3: list<Follow> followers
    4: i32 total
    5: bool has_more
}

// 获取关注数量请求
struct GetFollowCountRequest {
    1: string user_id
}

struct GetFollowCountResponse {
    1: i32 code
    2: string message
    3: i32 count
}

// 获取粉丝数量请求
struct GetFollowerCountRequest {
    1: string user_id
}

struct GetFollowerCountResponse {
    1: i32 code
    2: string message
    3: i32 count
}

// 获取共同关注请求
struct GetMutualFollowsRequest {
    1: string user_id
    2: string target_user_id
    3: i32 page
    4: i32 page_size
}

struct GetMutualFollowsResponse {
    1: i32 code
    2: string message
    3: list<Follow> mutual_follows
    4: i32 total
    5: bool has_more
}



service FollowService {
    FollowResponse Follow(1: FollowRequest req)
    UnfollowResponse Unfollow(1: UnfollowRequest req)  // 修改这里
    FollowResponse IsFollowing(1: FollowRequest req)
    GetFollowListResponse GetFollowList(1: GetFollowListRequest req)
    GetFollowerListResponse GetFollowerList(1: GetFollowerListRequest req)  // 添加这个方法
    CheckFollowStatusResponse CheckFollowStatus(1: CheckFollowStatusRequest req)
    GetFollowCountResponse GetFollowCount(1: GetFollowCountRequest req)
    GetFollowerCountResponse GetFollowerCount(1: GetFollowerCountRequest req)
    GetMutualFollowsResponse GetMutualFollows(1: GetMutualFollowsRequest req)
}