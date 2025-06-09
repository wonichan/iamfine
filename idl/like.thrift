namespace go like

struct LikeRequest {
    1: string user_id
    2: string target_id
    3: string target_type // post, comment
}

struct LikeResponse {
    1: i32 code
    2: string message
}

// 添加 UnlikeRequest 结构体
struct UnlikeRequest {
    1: string user_id
    2: string target_id
    3: string target_type // post, comment
}

// 添加 UnlikeResponse 结构体
struct UnlikeResponse {
    1: i32 code
    2: string message
}

// 添加 GetLikeListRequest 结构体
struct GetLikeListRequest {
    1: string user_id
    2: i32 page
    3: i32 page_size
}

// 添加 GetLikeListResponse 结构体
struct GetLikeListResponse {
    1: i32 code
    2: string message
    3: list<Like> likes
    4: i32 total
}

// 添加 Like 结构体
struct Like {
    1: string user_id
    2: string target_id
    3: string target_type
    4: i64 created_at
}

service LikeService {
    LikeResponse Like(1: LikeRequest req)
    UnlikeResponse Unlike(1: UnlikeRequest req)  // 修改这里
    LikeResponse IsLiked(1: LikeRequest req)
    GetLikeListResponse GetLikeList(1: GetLikeListRequest req)  // 添加这个方法
}