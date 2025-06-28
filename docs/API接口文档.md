# 微信小程序API接口文档

## 概述

本文档描述了微信小程序后端API的所有接口定义，包括用户管理、帖子管理、评论管理、评分管理、话题管理、轮播图管理、文件上传和统计等模块。

## 基础配置

- **基础URL**: `https://api.yourapp.com`
- **请求格式**: JSON
- **响应格式**: JSON
- **认证方式**: Bearer Token

## 通用响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 1. 用户管理模块

### 1.1 微信登录

**接口地址**: `POST /api/user/wx-login`

**请求参数**:
```json
{
  "code": "string" // 微信登录凭证
}
```

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "string",
    "user": {
      "id": "string",
      "openId": "string",
      "unionId": "string",
      "avatar": "string",
      "nickname": "string",
      "gender": 0,
      "city": "string",
      "province": "string",
      "country": "string",
      "tags": ["string"],
      "createdAt": "string",
      "updatedAt": "string"
    }
  }
}
```

### 1.2 获取用户信息

**接口地址**: `GET /api/user/info`

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
    "id": "string",
    "openId": "string",
    "unionId": "string",
    "avatar": "string",
    "nickname": "string",
    "gender": 0,
    "city": "string",
    "province": "string",
    "country": "string",
    "tags": ["string"],
    "createdAt": "string",
    "updatedAt": "string"
  }
}
```

### 1.3 更新用户信息

**接口地址**: `PUT /api/user/info`

**请求参数**:
```json
{
  "avatar": "string",
  "nickname": "string",
  "gender": 0,
  "city": "string",
  "province": "string",
  "country": "string",
  "tags": ["string"]
}
```

### 1.4 获取用户统计

**接口地址**: `GET /api/user/stats`

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "posts": 0,
    "comments": 0,
    "favorites": 0,
    "ratings": 0,
    "followers": 0,
    "following": 0
  }
}
```

### 1.5 获取未读消息数

**接口地址**: `GET /api/user/unread-count`

**响应数据**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "count": 0
  }
}
```

## 2. 帖子管理模块

### 2.1 获取帖子列表

**接口地址**: `GET /api/posts`

**请求参数**:
```
page: number (页码，默认1)
size: number (每页数量，默认10)
category: string (分类，可选)
tag: string (标签，可选)
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
        "userId": "string",
        "title": "string",
        "content": "string",
        "images": ["string"],
        "category": "string",
        "tags": ["string"],
        "likeCount": 0,
        "commentCount": 0,
        "favoriteCount": 0,
        "viewCount": 0,
        "score": 0,
        "ratingCount": 0,
        "isLiked": false,
        "isFavorited": false,
        "isRated": false,
        "userRating": 0,
        "author": {
          "id": "string",
          "nickname": "string",
          "avatar": "string"
        },
        "createdAt": "string",
        "updatedAt": "string"
      }
    ],
    "total": 0,
    "hasMore": true
  }
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

### 2.7 获取帖子详情

**接口地址**: `GET /api/posts/{id}`

**响应数据**: 单个帖子对象

### 2.8 创建帖子

**接口地址**: `POST /api/posts`

**请求参数**:
```json
{
  "title": "string",
  "content": "string",
  "images": ["string"],
  "category": "string",
  "tags": ["string"]
}
```

### 2.9 更新帖子

**接口地址**: `PUT /api/posts/{id}`

**请求参数**: 同创建帖子

### 2.10 删除帖子

**接口地址**: `DELETE /api/posts/{id}`

### 2.11 点赞帖子

**接口地址**: `POST /api/posts/{id}/like`

### 2.12 取消点赞

**接口地址**: `DELETE /api/posts/{id}/like`

### 2.13 收藏帖子

**接口地址**: `POST /api/posts/{id}/favorite`

### 2.14 取消收藏

**接口地址**: `DELETE /api/posts/{id}/favorite`

### 2.15 搜索帖子

**接口地址**: `GET /api/posts/search`

**请求参数**:
```
keyword: string (搜索关键词)
page: number
size: number
```

## 3. 评论管理模块

### 3.1 获取评论列表

**接口地址**: `GET /api/comments`

**请求参数**:
```
postId: string (帖子ID)
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
        "postId": "string",
        "userId": "string",
        "parentId": "string",
        "content": "string",
        "likeCount": 0,
        "isLiked": false,
        "author": {
          "id": "string",
          "nickname": "string",
          "avatar": "string"
        },
        "replies": [],
        "createdAt": "string",
        "updatedAt": "string"
      }
    ],
    "total": 0,
    "hasMore": true
  }
}
```

### 3.2 创建评论

**接口地址**: `POST /api/comments`

**请求参数**:
```json
{
  "postId": "string",
  "parentId": "string",
  "content": "string"
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

## 7. 文件上传模块

### 7.1 上传图片

**接口地址**: `POST /api/upload/image`

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
    "type": "string"
  }
}
```

### 7.2 批量上传图片

**接口地址**: `POST /api/upload/images`

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