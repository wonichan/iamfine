namespace go user

struct User {
    1: string id
    2: string username
    3: string nickname
    4: string avatar
    5: string phone
    6: string email
    7: i64 created_at
    8: i64 updated_at
}

struct RegisterRequest {
    1: string username
    2: string password
    3: string phone
    4: string code
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
    3: string token
    4: User user
}

struct GetUserRequest {
    1: string user_id
}

struct GetUserResponse {
    1: i32 code
    2: string message
    3: User user
}

service UserService {
    RegisterResponse Register(1: RegisterRequest req)
    LoginResponse Login(1: LoginRequest req)
    GetUserResponse GetUser(1: GetUserRequest req)
}