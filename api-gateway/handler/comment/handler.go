package comment

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// GetCommentListHandler 获取评论列表
func GetCommentListHandler(ctx context.Context, c *app.RequestContext) {
	GetCommentList(ctx, c)
}

// CreateCommentHandler 创建评论
func CreateCommentHandler(ctx context.Context, c *app.RequestContext) {
	CreateComment(ctx, c)
}

// DeleteCommentHandler 删除评论
func DeleteCommentHandler(ctx context.Context, c *app.RequestContext) {
	DeleteComment(ctx, c)
}

// LikeCommentHandler 点赞评论
func LikeCommentHandler(ctx context.Context, c *app.RequestContext) {
	LikeComment(ctx, c)
}

// UnlikeCommentHandler 取消点赞评论
func UnlikeCommentHandler(ctx context.Context, c *app.RequestContext) {
	UnlikeComment(ctx, c)
}