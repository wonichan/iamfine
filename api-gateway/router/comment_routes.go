package router

import (
	"hupu/api-gateway/handler/comment"
	"hupu/api-gateway/handler/like"
	"hupu/api-gateway/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterCommentRoutes 注册评论相关路由
func RegisterCommentRoutes(h *server.Hertz) {
	// API版本组
	apiV1 := h.Group("/api")

	// 评论相关路由（需要认证）
	authGroup := apiV1.Group("", middleware.AuthMiddleware())
	{
		// 评论管理
		authGroup.GET("/comments", comment.GetCommentListHandler)         // 获取评论列表
		authGroup.POST("/comments", comment.CreateCommentHandler)          // 创建评论
		authGroup.DELETE("/comments/:id", comment.DeleteCommentHandler)   // 删除评论
		// 评论点赞
		authGroup.POST("/comments/:id/like", like.LikeCommentHandler)   // 点赞评论
		authGroup.DELETE("/comments/:id/like", like.UnlikeCommentHandler) // 取消点赞评论
	}

	// 无需认证的评论点赞统计路由
	apiV1.GET("/comments/:id/like/count", like.GetCommentLikeCountHandler)
	apiV1.GET("/comments/:id/like/status", like.CheckCommentLikeStatusHandler)
}