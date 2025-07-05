# 虎扑情感社区API接口文档
## 概述

本文档描述了虎扑情感社区后端API的所有接口定义，包括用户管理、帖子管理、评论管理、点赞管理、关注管理、通知管理等模块。基于微服务架构，使用Thrift作为RPC协议，Hertz作为HTTP框架。

## 基础配置

- **基础URL**: `http://localhost:8080`
- **API版本**: `/api/v1`
- **请求格式**: JSON
- **响应格式**: JSON
- **认证方式**: Bearer Token (JWT)
- **架构**: 微服务架构 (API Gateway + 多个微服务)

## 通用响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

## 7. 错误码说明

### 7.1 HTTP状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未授权，需要登录 |
| 403 | 禁止访问，权限不足 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 422 | 请求参数验证失败 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |
| 502 | 网关错误 |
| 503 | 服务不可用 |

### 7.2 业务错误码

| 错误码 | 说明 |
|--------|------|
| 10001 | 用户名已存在 |
| 10002 | 用户名或密码错误 |
| 10003 | 用户不存在 |
| 10004 | 用户已被禁用 |
| 10005 | 微信授权失败 |
| 20001 | 帖子不存在 |
| 20002 | 帖子已被删除 |
| 20003 | 无权限操作此帖子 |
| 20004 | 帖子内容违规 |
| 30001 | 评论不存在 |
| 30002 | 评论已被删除 |
| 30003 | 无权限操作此评论 |
| 30004 | 评论内容违规 |
| 40001 | 话题不存在 |
| 40002 | 话题名称已存在 |
| 50001 | 文件格式不支持 |
| 50002 | 文件大小超出限制 |
| 50003 | 文件上传失败 |
| 60001 | 已关注该用户 |
| 60002 | 未关注该用户 |
| 60003 | 不能关注自己 |
| 70001 | 已点赞 |
| 70002 | 未点赞 |
| 80001 | 已收藏 |
| 80002 | 未收藏 |
| 90001 | 评分超出范围 |
| 90002 | 已评分过 |
| 90003 | 不能给自己的帖子评分 |

## 1. 用户管理模块

### 1.1 用户注册

**接口地址**: `POST /api/v1/user/register`

**请求参数**:
```json
{
  "username": "string",
  "password": "string",
  "nickname": "string",
  "phone": "string",
  "email": "string"
}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "user": {
    "id": "string",
    "username": "string",
    "nickname": "string",
    "avatar": "string",
    "phone": "string",
    "email": "string",
    "status": 1,
    "bio": "string",
    "relationship_status": 1,
    "age_group": 1,
    "location": "string",
    "post_count": 0,
    "comment_count": 0,
    "like_count": 0,
    "collect_count": 0,
    "average_score": 0.0,
    "follower_count": 0,
    "following_count": 0,
    "is_verified": false,
    "tags": ["string"],
    "created_at": 1640995200,
    "updated_at": 1640995200
  }
}
```

### 1.2 用户登录

**接口地址**: `POST /api/v1/user/login`

**请求参数**:
```json
{
  "username": "string",
  "password": "string"
}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "user": {
    "id": "string",
    "username": "string",
    "nickname": "string",
    "avatar": "string",
    "phone": "string",
    "email": "string",
    "status": 1,
    "bio": "string",
    "relationship_status": 1,
    "age_group": 1,
    "location": "string",
    "post_count": 0,
    "comment_count": 0,
    "like_count": 0,
    "collect_count": 0,
    "average_score": 0.0,
    "follower_count": 0,
    "following_count": 0,
    "is_verified": false,
    "tags": ["string"],
    "created_at": 1640995200,
    "updated_at": 1640995200
  }
}
```

### 1.3 微信登录

**接口地址**: `POST /api/v1/user/wx-login`

**请求参数**:
```json
{
  "code": "string"
}
```

**响应数据**: 同用户登录

### 1.4 获取用户信息

**接口地址**: `GET /api/v1/user/info`

**请求头**:
```
Authorization: Bearer {token}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "user": {
    "id": "string",
    "username": "string",
    "nickname": "string",
    "avatar": "string",
    "phone": "string",
    "email": "string",
    "status": 1,
    "bio": "string",
    "relationship_status": 1,
    "age_group": 1,
    "location": "string",
    "post_count": 0,
    "comment_count": 0,
    "like_count": 0,
    "collect_count": 0,
    "average_score": 0.0,
    "follower_count": 0,
    "following_count": 0,
    "is_verified": false,
    "tags": ["string"],
    "anonymous_profiles": [
      {
        "id": "string",
        "user_id": "string",
        "anonymous_name": "string",
        "avatar_color": "string",
        "is_active": true,
        "created_at": 1640995200
      }
    ],
    "created_at": 1640995200,
    "updated_at": 1640995200
  }
}
```

### 1.5 获取用户资料

**接口地址**: `GET /api/v1/user/profile/{id}`

**路径参数**:
- `id`: 用户ID

**响应数据**: 同获取用户信息

### 1.6 更新用户信息

**接口地址**: `PUT /api/v1/user/{id}`

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "nickname": "string",
  "avatar": "string",
  "phone": "string",
  "email": "string",
  "bio": "string",
  "relationship_status": 1,
  "age_group": 1,
  "location": "string",
  "tags": ["string"]
}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "user": {
    // 用户完整信息
  }
}
```

### 1.7 获取用户统计

**接口地址**: `GET /api/v1/user/stats`

**请求头**:
```
Authorization: Bearer {token}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "post_count": 0,
    "comment_count": 0,
    "like_count": 0,
    "collect_count": 0,
    "average_score": 0.0,
    "follower_count": 0,
    "following_count": 0
  }
}
```

### 1.8 获取未读消息数

**接口地址**: `GET /api/v1/user/unread-count`

**请求头**:
```
Authorization: Bearer {token}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "unread_count": 0
}
```

### 1.9 关注用户

**接口地址**: `POST /api/v1/user/{id}/follow`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 要关注的用户ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success"
}
```

### 1.10 取消关注用户

**接口地址**: `DELETE /api/v1/user/{id}/follow`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 要取消关注的用户ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success"
}
```

### 1.11 获取用户粉丝列表

**接口地址**: `GET /api/v1/user/{id}/followers`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 用户ID

**请求参数**:
```
page: number (页码，默认1)
page_size: number (每页数量，默认10)
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "users": [
    {
      "id": "string",
      "nickname": "string",
      "avatar": "string"
    }
  ],
  "total": 0
}
```

### 1.12 获取用户关注列表

**接口地址**: `GET /api/v1/user/{id}/following`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 用户ID

**请求参数**:
```
page: number (页码，默认1)
page_size: number (每页数量，默认10)
```

**响应数据**: 同获取用户粉丝列表

## 2. 帖子管理模块

### 2.1 获取帖子列表

**接口地址**: `GET /api/v1/posts`

**请求参数**:
```
page: number (页码，默认1)
page_size: number (每页数量，默认10)
user_id: string (用户ID，可选)
topic_id: string (话题ID，可选)
category: number (分类，可选: 1-日常分享, 2-恋爱日常, 3-婚姻围城, 4-家庭关系, 5-情感求助, 6-我要吐槽, 99-其他)
sort_type: string (排序类型，可选: latest, hot, score)
is_anonymous: boolean (是否只看匿名帖子，可选)
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "posts": [
    {
      "id": "string",
      "user_id": "string",
      "title": "string",
      "content": "string",
      "images": ["string"],
      "like_count": 0,
      "comment_count": 0,
      "score": 0.0,
      "created_at": 1640995200,
      "updated_at": 1640995200,
      "topic_id": "string",
      "category": 1,
      "is_anonymous": false,
      "anonymous_name": "string",
      "view_count": 0,
      "share_count": 0,
      "collect_count": 0,
      "is_hot": false,
      "is_top": false,
      "location": "string",
      "tags": ["string"]
    }
  ],
  "total": 0
}
```

### 2.2 获取推荐帖子

**接口地址**: `GET /api/posts/recommend`

**请求参数**: 同获取帖子列表

### 2.3 获取热门帖子

**接口地址**: `GET /api/posts/hot`

**请求参数**: 同获取帖子列表

### 2.4 获取高分帖子

**接口地址**: `GET /api/posts/high-score`

**请求参数**: 同获取帖子列表

### 2.5 获取低分帖子

**接口地址**: `GET /api/posts/low-score`

**请求参数**: 同获取帖子列表

### 2.6 获取争议帖子

**接口地址**: `GET /api/posts/controversial`

**请求参数**: 同获取帖子列表

### 2.2 获取帖子详情

**接口地址**: `GET /api/v1/posts/{id}`

**路径参数**:
- `id`: 帖子ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "post": {
    "id": "string",
    "user_id": "string",
    "title": "string",
    "content": "string",
    "images": ["string"],
    "like_count": 0,
    "comment_count": 0,
    "score": 0.0,
    "created_at": 1640995200,
    "updated_at": 1640995200,
    "topic_id": "string",
    "category": 1,
    "is_anonymous": false,
    "anonymous_name": "string",
    "view_count": 0,
    "share_count": 0,
    "collect_count": 0,
    "is_hot": false,
    "is_top": false,
    "location": "string",
    "tags": ["string"]
  }
}
```

### 2.3 创建帖子

**接口地址**: `POST /api/v1/posts`

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "title": "string",
  "content": "string",
  "images": ["string"],
  "topic_id": "string",
  "category": 1,
  "is_anonymous": false,
  "location": "string",
  "tags": ["string"]
}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "post": {
    // 帖子完整信息
  }
}
```

### 2.4 更新帖子

**接口地址**: `PUT /api/v1/posts/{id}`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 帖子ID

**请求参数**: 同创建帖子

**响应数据**: 同创建帖子

### 2.5 删除帖子

**接口地址**: `DELETE /api/v1/posts/{id}`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 帖子ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success"
}
```

### 2.6 获取推荐帖子

**接口地址**: `GET /api/v1/posts/recommend`

**请求参数**: 同获取帖子列表

### 2.7 获取热门帖子

**接口地址**: `GET /api/v1/posts/hot`

**请求参数**: 同获取帖子列表

### 2.8 获取高分帖子

**接口地址**: `GET /api/v1/posts/high-score`

**请求参数**: 同获取帖子列表

### 2.9 获取低分帖子

**接口地址**: `GET /api/v1/posts/low-score`

**请求参数**: 同获取帖子列表

### 2.10 获取争议帖子

**接口地址**: `GET /api/v1/posts/controversial`

**请求参数**: 同获取帖子列表

### 2.11 搜索帖子

**接口地址**: `GET /api/v1/posts/search`

**请求参数**:
```
keyword: string (搜索关键词)
page: number (页码，默认1)
page_size: number (每页数量，默认10)
category: number (分类，可选)
sort_type: string (排序类型，可选)
```

**响应数据**: 同获取帖子列表

### 2.12 点赞帖子

**接口地址**: `POST /api/v1/posts/{id}/like`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 帖子ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success"
}
```

### 2.13 取消点赞帖子

**接口地址**: `DELETE /api/v1/posts/{id}/like`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 帖子ID

**响应数据**: 同点赞帖子

### 2.14 收藏帖子

**接口地址**: `POST /api/v1/posts/{id}/collect`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 帖子ID

**响应数据**: 同点赞帖子

### 2.15 取消收藏帖子

**接口地址**: `DELETE /api/v1/posts/{id}/collect`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 帖子ID

**响应数据**: 同点赞帖子

### 2.16 获取收藏的帖子

**接口地址**: `GET /api/v1/posts/collected`

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:
```
page: number (页码，默认1)
page_size: number (每页数量，默认10)
```

**响应数据**: 同获取帖子列表

### 2.17 评分帖子

**接口地址**: `POST /api/v1/posts/{id}/rate`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 帖子ID

**请求参数**:
```json
{
  "score": 8.5,
  "comment": "string"
}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "average_score": 8.2,
  "total_ratings": 15
}
```

### 2.18 获取用户对帖子的评分

**接口地址**: `GET /api/v1/posts/{id}/rating`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 帖子ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "score": 8.5,
  "comment": "string"
}
```

### 2.19 更新帖子评分

**接口地址**: `PUT /api/v1/posts/{id}/rating`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 帖子ID

**请求参数**: 同评分帖子

**响应数据**: 同评分帖子

### 2.20 删除帖子评分

**接口地址**: `DELETE /api/v1/posts/{id}/rating`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 帖子ID

**响应数据**: 同评分帖子

### 2.21 获取评分排行榜

**接口地址**: `GET /api/v1/posts/rating/rank`

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:
```
page: number (页码，默认1)
page_size: number (每页数量，默认10)
rank_type: string (排行类型: daily_high, daily_low, weekly_best, controversial)
date: string (指定日期，格式: 2024-01-01，可选)
```

**响应数据**: 同获取帖子列表

### 2.22 获取帖子点赞数

**接口地址**: `GET /api/v1/posts/{id}/like/count`

**路径参数**:
- `id`: 帖子ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "count": 25
}
```

### 2.23 检查帖子点赞状态

**接口地址**: `GET /api/v1/posts/{id}/like/status`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 帖子ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "is_liked": true
}
```

## 3. 话题管理模块

### 3.1 获取话题列表

**接口地址**: `GET /api/v1/posts/topics`

**请求参数**:
```
page: number (页码，默认1)
page_size: number (每页数量，默认10)
sort_type: string (排序类型，可选: hot, latest, participant)
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "topics": [
    {
      "id": "string",
      "name": "string",
      "description": "string",
      "icon": "string",
      "color": "string",
      "participant_count": 0,
      "created_at": 1640995200,
      "updated_at": 1640995200
    }
  ],
  "total": 0
}
```

### 3.2 获取话题详情

**接口地址**: `GET /api/v1/posts/topics/{id}`

**路径参数**:
- `id`: 话题ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "topic": {
    "id": "string",
    "name": "string",
    "description": "string",
    "icon": "string",
    "color": "string",
    "participant_count": 0,
    "created_at": 1640995200,
    "updated_at": 1640995200
  }
}
```

### 3.3 获取热门话题

**接口地址**: `GET /api/v1/posts/topics/hot`

**请求参数**: 同获取话题列表

### 3.4 获取话题分类

**接口地址**: `GET /api/v1/posts/topics/categories`

**响应数据**: 同获取话题列表

### 3.5 搜索话题

**接口地址**: `GET /api/v1/posts/topics/search`

**请求参数**:
```
keyword: string (搜索关键词)
page: number (页码，默认1)
page_size: number (每页数量，默认10)
```

**响应数据**: 同获取话题列表

### 3.6 创建话题

**接口地址**: `POST /api/v1/posts/topics`

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "name": "string",
  "description": "string",
  "icon": "string",
  "color": "string"
}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "topic": {
    // 话题完整信息
  }
}
```

### 3.7 分享话题

**接口地址**: `POST /api/v1/posts/topics/{id}/share`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 话题ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success"
}
```

## 4. 评论管理模块

### 4.1 获取评论列表

**接口地址**: `GET /api/comments`

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:
```
post_id: string (帖子ID)
page: number (页码，默认1)
page_size: number (每页数量，默认10)
sort_type: string (排序类型，可选: latest, hot, oldest, score_high, score_low)
parent_id: string (获取特定评论的回复，可选)
include_replies: boolean (是否包含回复，默认true)
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": "string",
        "post_id": "string",
        "user_id": "string",
        "content": "string",
        "created_at": 1640995200,
        "updated_at": 1640995200,
        "parent_id": "string",
        "like_count": 0,
        "reply_count": 0,
        "is_anonymous": false,
        "anonymous_name": "string",
        "anonymous_color": "string",
        "images": ["string"],
        "is_deleted": false,
        "location": "string",
        "replies": [],
        "is_liked": false,
        "author": {
          "id": "string",
          "nickname": "string",
          "avatar": "string"
        },
        "score": 0.0,
        "rating_count": 0,
        "is_rated": false
      }
    ],
    "total": 0,
    "hasMore": true
  }
}
```

### 4.2 创建评论

**接口地址**: `POST /api/comments`

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "post_id": "string",
  "content": "string",
  "parent_id": "string",
  "is_anonymous": false,
  "anonymous_profile_id": "string",
  "images": ["string"],
  "location": "string"
}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "comment": {
    // 评论完整信息
  }
}
```

### 4.3 删除评论

**接口地址**: `DELETE /api/comments/{id}`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 评论ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success"
}
```

### 4.4 点赞评论

**接口地址**: `POST /api/comments/{id}/like`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 评论ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success"
}
```

### 4.5 取消点赞评论

**接口地址**: `DELETE /api/comments/{id}/like`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 评论ID

**响应数据**: 同点赞评论

### 4.6 获取评论点赞数

**接口地址**: `GET /api/comments/{id}/like/count`

**路径参数**:
- `id`: 评论ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "count": 15
}
```

### 4.7 检查评论点赞状态

**接口地址**: `GET /api/comments/{id}/like/status`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 评论ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "is_liked": true
}
```

### 3.3 删除评论

**接口地址**: `DELETE /api/comments/{id}`

### 3.4 点赞评论

**接口地址**: `POST /api/comments/{id}/like`

### 3.5 取消点赞评论

**接口地址**: `DELETE /api/comments/{id}/like`

## 4. 评分管理模块

### 4.1 对帖子评分

**接口地址**: `POST /api/ratings`

**请求参数**:
```json
{
  "postId": "string",
  "score": 5
}
```

### 4.2 获取用户对帖子的评分

**接口地址**: `GET /api/ratings/{postId}`

### 4.3 更新评分

**接口地址**: `PUT /api/ratings/{postId}`

**请求参数**:
```json
{
  "score": 5
}
```

### 4.4 删除评分

**接口地址**: `DELETE /api/ratings/{postId}`

## 5. 话题管理模块

### 5.1 获取热门话题

**接口地址**: `GET /api/topics/hot`

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "string",
      "title": "string",
      "description": "string",
      "emoji": "string",
      "postCount": 0,
      "todayCount": 0,
      "backgroundColor": "string",
      "iconBackground": "string",
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}
```

### 5.2 获取话题分类

**接口地址**: `GET /api/topics/categories`

**响应数据**: 同热门话题

### 5.3 搜索话题

**接口地址**: `GET /api/topics/search`

**请求参数**:
```
keyword: string
page: number
size: number
```

### 5.4 创建话题

**接口地址**: `POST /api/topics`

**请求参数**:
```json
{
  "title": "string",
  "description": "string",
  "emoji": "string",
  "backgroundColor": "string",
  "iconBackground": "string"
}
```

### 5.5 获取话题详情

**接口地址**: `GET /api/topics/{id}`

### 5.6 分享话题

**接口地址**: `POST /api/topics/{id}/share`

## 6. 轮播图管理模块

### 6.1 获取轮播图列表

**接口地址**: `GET /api/carousels`

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": "string",
      "title": "string",
      "image": "string",
      "link": "string",
      "linkType": "string",
      "sort": 0,
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}
```

## 5. 通知管理模块

### 5.1 获取通知列表

**接口地址**: `GET /api/v1/notifications`

**请求头**:
```
Authorization: Bearer {token}
```

**请求参数**:
```
page: number (页码，默认1)
page_size: number (每页数量，默认10)
type: string (通知类型，可选: like, comment, follow, system)
is_read: boolean (是否已读，可选)
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": "string",
        "user_id": "string",
        "type": "string",
        "title": "string",
        "content": "string",
        "data": {},
        "is_read": false,
        "created_at": 1640995200,
        "updated_at": 1640995200,
        "from_user": {
          "id": "string",
          "nickname": "string",
          "avatar": "string"
        }
      }
    ],
    "total": 0,
    "unread_count": 5
  }
}
```

### 5.2 标记通知为已读

**接口地址**: `PUT /api/v1/notifications/{id}/read`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 通知ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success"
}
```

### 5.3 标记所有通知为已读

**接口地址**: `PUT /api/v1/notifications/read-all`

**请求头**:
```
Authorization: Bearer {token}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success"
}
```

### 5.4 删除通知

**接口地址**: `DELETE /api/v1/notifications/{id}`

**请求头**:
```
Authorization: Bearer {token}
```

**路径参数**:
- `id`: 通知ID

**响应数据**:
```json
{
  "code": 200,
  "message": "success"
}
```

### 5.5 获取未读通知数量

**接口地址**: `GET /api/v1/notifications/unread-count`

**请求头**:
```
Authorization: Bearer {token}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "count": 5
}
```

## 6. 文件上传模块

### 6.1 上传图片

**接口地址**: `POST /api/v1/upload/image`

**请求头**:
```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**请求参数**: FormData (multipart/form-data)
```
file: File (图片文件)
type: string (上传类型: avatar, post, comment等)
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "url": "string",
    "filename": "string",
    "size": 0,
    "width": 1920,
    "height": 1080,
    "format": "jpg",
    "type": "string"
  }
}
```

### 6.2 批量上传图片

**接口地址**: `POST /api/v1/upload/images`

**请求头**:
```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**请求参数**: FormData
```
files: File[] (多个图片文件)
type: string
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "url": "string",
      "filename": "string",
      "size": 0,
      "width": 1920,
      "height": 1080,
      "format": "jpg",
      "type": "string"
    }
  ]
}
```

## 8. 统计模块

### 8.1 获取应用统计

**接口地址**: `GET /api/stats/app`

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "totalUsers": 0,
    "totalPosts": 0,
    "totalComments": 0,
    "totalTopics": 0,
    "todayUsers": 0,
    "todayPosts": 0,
    "todayComments": 0
  }
}
```

### 8.2 获取用户统计

**接口地址**: `GET /api/stats/user`

**响应数据**: 参考用户统计接口

## 9. 消息通知模块

### 9.1 获取消息列表

**接口地址**: `GET /api/messages`

**请求参数**:
```
type: string (消息类型: like, comment, follow等)
page: number
size: number
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": "string",
        "type": "string",
        "title": "string",
        "content": "string",
        "data": {},
        "isRead": false,
        "createdAt": "string"
      }
    ],
    "total": 0,
    "hasMore": true
  }
}
```

### 9.2 标记消息已读

**接口地址**: `PUT /api/messages/{id}/read`

### 9.3 批量标记已读

**接口地址**: `PUT /api/messages/read-all`

## 10. 搜索模块

### 10.1 综合搜索

**接口地址**: `GET /api/search`

**请求参数**:
```
keyword: string
type: string (all, post, topic, user)
page: number
size: number
```

### 10.2 获取搜索建议

**接口地址**: `GET /api/search/suggestions`

**请求参数**:
```
keyword: string
limit: number (默认10)
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": ["string"]
}
```

### 10.3 获取热门搜索

**接口地址**: `GET /api/search/hot`

**响应数据**: 同搜索建议

## 注意事项

1. 所有需要用户身份验证的接口都需要在请求头中携带 `Authorization: Bearer {token}`
2. 分页参数中，`page` 从1开始，`size` 建议不超过50
3. 图片上传支持的格式：jpg, jpeg, png, gif，单个文件大小不超过5MB
4. 所有时间字段均为ISO 8601格式的UTC时间
5. 评分范围为1-10分
6. 搜索关键词长度不超过50个字符

## SDK使用说明

本项目使用微信小程序原生的 `wx.request` API进行网络请求，并封装了统一的API服务类 `apiService`，提供了以下功能：

- 统一的请求拦截和响应处理
- 自动的Token管理
- 错误处理和用户提示
- 请求参数格式化
- 文件上传支持

使用示例：
```typescript
import { PostService } from './services/index'

// 获取帖子列表
const posts = await PostService.getPostList({ page: 1, size: 10 })

// 创建帖子
const newPost = await PostService.createPost({
  title: '标题',
  content: '内容',
  images: ['image1.jpg'],
  category: '情感',
  tags: ['恋爱']
})
```

所有服务类都已在 `services/index.ts` 中定义并导出，可直接在页面中引入使用。