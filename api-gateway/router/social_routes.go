package router

import (
	"hupu/api-gateway/handler"
	"hupu/api-gateway/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterSocialRoutes 注册社交相关路由（评论、点赞、关注、通知）
func RegisterSocialRoutes(h *server.Hertz) {
	// API版本组
	apiV1 := h.Group("/api/v1")

	// 需要认证的路由
	authGroup := apiV1.Group("", middleware.AuthMiddleware())
	{
		// 评论相关
		authGroup.POST("/comment", handler.CreateComment)
		authGroup.GET("/comments/:post_id", handler.GetCommentList)

		// 点赞相关
		authGroup.POST("/like", handler.Like)
		authGroup.DELETE("/like", handler.Unlike)
		authGroup.GET("/likes", handler.GetLikeList)

		// 关注相关
		authGroup.POST("/follow", handler.Follow)
		authGroup.DELETE("/follow", handler.Unfollow)
		authGroup.GET("/follows/:user_id", handler.GetFollowList)
		authGroup.GET("/followers/:user_id", handler.GetFollowerList)

		// 通知相关
		authGroup.GET("/notifications", handler.GetNotificationList)
	}
}