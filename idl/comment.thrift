namespace go comment

struct Comment {
    1: string id
    2: string post_id
    3: string user_id
    4: optional string parent_id
    5: string content
    6: i32 likeCount
    7: i64 created_at
    8: i64 updated_at
}

struct CreateCommentRequest {
    1: string post_id
    2: string user_id
    3: optional string parent_id
    4: string content
}

struct CreateCommentResponse {
    1: i32 code
    2: string message
    3: Comment comment
}

struct GetCommentListRequest {
    1: string post_id
    2: optional string parent_id
    3: i32 page
    4: i32 page_size
}

struct GetCommentListResponse {
    1: i32 code
    2: string message
    3: list<Comment> comments
    4: i32 total
}

service CommentService {
    CreateCommentResponse CreateComment(1: CreateCommentRequest req)
    GetCommentListResponse GetCommentList(1: GetCommentListRequest req)
}