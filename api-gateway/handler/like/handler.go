package like

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// ==================== 帖子点赞Handler ====================

// LikePostHandler 点赞帖子
func LikePostHandler(ctx context.Context, c *app.RequestContext) {
	LikePost(ctx, c)
}

// UnlikePostHandler 取消点赞帖子
func UnlikePostHandler(ctx context.Context, c *app.RequestContext) {
	UnlikePost(ctx, c)
}

// GetPostLikeCountHandler 获取帖子点赞数量
func GetPostLikeCountHandler(ctx context.Context, c *app.RequestContext) {
	GetPostLikeCount(ctx, c)
}

// CheckPostLikeStatusHandler 检查帖子点赞状态
func CheckPostLikeStatusHandler(ctx context.Context, c *app.RequestContext) {
	CheckPostLikeStatus(ctx, c)
}

// ==================== 评论点赞Handler ====================

// LikeCommentHandler 点赞评论
func LikeCommentHandler(ctx context.Context, c *app.RequestContext) {
	LikeComment(ctx, c)
}

// UnlikeCommentHandler 取消点赞评论
func UnlikeCommentHandler(ctx context.Context, c *app.RequestContext) {
	UnlikeComment(ctx, c)
}

// GetCommentLikeCountHandler 获取评论点赞数量
func GetCommentLikeCountHandler(ctx context.Context, c *app.RequestContext) {
	GetCommentLikeCount(ctx, c)
}

// CheckCommentLikeStatusHandler 检查评论点赞状态
func CheckCommentLikeStatusHandler(ctx context.Context, c *app.RequestContext) {
	CheckCommentLikeStatus(ctx, c)
}

// ==================== 通用点赞Handler ====================

// GetLikeListHandler 获取用户点赞列表
func GetLikeListHandler(ctx context.Context, c *app.RequestContext) {
	GetLikeList(ctx, c)
}

// IsLikedHandler 检查是否已点赞
func IsLikedHandler(ctx context.Context, c *app.RequestContext) {
	IsLiked(ctx, c)
}

// GetLikeCountHandler 获取点赞数量
func GetLikeCountHandler(ctx context.Context, c *app.RequestContext) {
	GetLikeCount(ctx, c)
}

// GetLikeUsersHandler 获取点赞用户列表
func GetLikeUsersHandler(ctx context.Context, c *app.RequestContext) {
	GetLikeUsers(ctx, c)
}

// ==================== 兼容性Handler ====================

// LikeHandler 兼容旧的点赞接口
func LikeHandler(ctx context.Context, c *app.RequestContext) {
	Like(ctx, c)
}

// UnlikeHandler 兼容旧的取消点赞接口
func UnlikeHandler(ctx context.Context, c *app.RequestContext) {
	Unlike(ctx, c)
}