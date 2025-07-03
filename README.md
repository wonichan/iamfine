# 仿虎扑微信小程序后端系统

## 项目简介

本项目是一个仿虎扑的微信小程序后端系统，采用微服务架构设计，使用Go语言开发。系统包含用户管理、帖子管理、评论系统、点赞功能、关注系统和通知服务等核心功能。

## 技术栈

- **Web框架**: Hertz v0.10.0
- **RPC框架**: Kitex v0.13.1
- **数据库**: MySQL 8.0
- **缓存**: Redis v9.7.0
- **消息队列**: Kafka
- **对象存储**: 阿里云OSS
- **分布式ID**: xid v1.6.0
- **容器化**: Docker & Docker Compose

## 系统架构

### 服务列表

1. **API网关** (端口: 8080) - 统一入口，负责路由转发和认证
2. **用户服务** (端口: 8001) - 用户注册、登录、信息管理
3. **帖子服务** (端口: 8002) - 帖子发布、查看、管理
4. **评论服务** (端口: 8003) - 评论发布、回复功能
5. **点赞服务** (端口: 8004) - 点赞、取消点赞
6. **关注服务** (端口: 8005) - 关注、取消关注、粉丝管理
7. **通知服务** (端口: 8006) - 消息推送和通知

### 数据库设计

#### 用户表 (users)
- id: 用户ID (主键)
- username: 用户名 (唯一)
- password: 密码 (加密)
- nickname: 昵称
- avatar: 头像URL
- phone: 手机号 (唯一)
- email: 邮箱
- created_at: 创建时间
- updated_at: 更新时间

#### 帖子表 (posts)
- id: 帖子ID (主键)
- user_id: 用户ID (外键)
- title: 标题
- content: 内容
- images: 图片列表 (JSON)
- like_count: 点赞数
- comment_count: 评论数
- created_at: 创建时间
- updated_at: 更新时间

#### 评论表 (comments)
- id: 评论ID (主键)
- post_id: 帖子ID (外键)
- user_id: 用户ID (外键)
- parent_id: 父评论ID (可选)
- content: 评论内容
- like_count: 点赞数
- created_at: 创建时间
- updated_at: 更新时间

#### 点赞表 (likes)
- id: 点赞ID (主键)
- user_id: 用户ID (外键)
- target_id: 目标ID (帖子或评论)
- target_type: 目标类型 (1:帖子, 2:评论)
- created_at: 创建时间

#### 关注表 (follows)
- id: 关注ID (主键)
- follower_id: 关注者ID (外键)
- following_id: 被关注者ID (外键)
- created_at: 创建时间

#### 通知表 (notifications)
- id: 通知ID (主键)
- user_id: 用户ID (外键)
- title: 通知标题
- content: 通知内容
- type: 通知类型 (1:点赞, 2:评论, 3:关注)
- target_id: 目标ID
- is_read: 是否已读
- created_at: 创建时间

## 快速开始

### 环境要求

- Go 1.21+
- Docker & Docker Compose
- MySQL 8.0
- Redis 7.0
- Kafka

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd hupu