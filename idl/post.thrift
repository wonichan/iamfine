namespace go post

struct Post {
    1: string id
    2: string user_id
    3: string title
    4: string content
    5: list<string> images
    6: i32 like_count
    7: i32 comment_count
    8: i64 created_at
    9: i64 updated_at
}

struct CreatePostRequest {
    1: string user_id
    2: string title
    3: string content
    4: list<string> images
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
    1: i32 page
    2: i32 page_size
    3: optional string user_id
}

struct GetPostListResponse {
    1: i32 code
    2: string message
    3: list<Post> posts
    4: i32 total
}

service PostService {
    CreatePostResponse CreatePost(1: CreatePostRequest req)
    GetPostResponse GetPost(1: GetPostRequest req)
    GetPostListResponse GetPostList(1: GetPostListRequest req)
}