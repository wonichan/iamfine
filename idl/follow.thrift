namespace go follow

struct FollowRequest {
    1: string follower_id
    2: string following_id
}

struct FollowResponse {
    1: i32 code
    2: string message
    3: bool is_following
}

// 添加 UnfollowRequest 结构体
struct UnfollowRequest {
    1: string follower_id
    2: string following_id
}

// 添加 UnfollowResponse 结构体
struct UnfollowResponse {
    1: i32 code
    2: string message
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
    3: list<string> user_ids
    4: i32 total
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
}

// 添加 Follow 结构体
struct Follow {
    1: string follower_id
    2: string following_id
    3: i64 created_at
}

service FollowService {
    FollowResponse Follow(1: FollowRequest req)
    UnfollowResponse Unfollow(1: UnfollowRequest req)  // 修改这里
    FollowResponse IsFollowing(1: FollowRequest req)
    GetFollowListResponse GetFollowList(1: GetFollowListRequest req)
    GetFollowerListResponse GetFollowerList(1: GetFollowerListRequest req)  // 添加这个方法
}