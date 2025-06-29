package router

import (
	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/follow"
	"hupu/api-gateway/handler/like"
	"hupu/api-gateway/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterSocialRoutes 注册社交相关路由（点赞、关注、通知）
func RegisterSocialRoutes(h *server.Hertz) {
	// API版本组
	apiV1 := h.Group("/api/v1")

	// 需要认证的路由
	authGroup := apiV1.Group("", middleware.AuthMiddleware())
	{
		// 点赞相关（兼容旧接口）
		authGroup.POST("/like", like.LikeHandler)
		authGroup.DELETE("/like", like.UnlikeHandler)
		authGroup.GET("/likes", like.GetLikeListHandler)
		// 点赞查询接口
		authGroup.GET("/like/status", like.IsLikedHandler)
		authGroup.GET("/like/count", like.GetLikeCountHandler)
		authGroup.GET("/like/users", like.GetLikeUsersHandler)

		// 关注相关
		authGroup.POST("/follow", follow.FollowHandler)
		authGroup.DELETE("/follow", follow.UnfollowHandler)
		authGroup.GET("/follows/:user_id", follow.GetFollowListHandler)
		authGroup.GET("/followers/:user_id", follow.GetFollowerListHandler)
		authGroup.GET("/follow/status", follow.CheckFollowStatusHandler)
		authGroup.GET("/follow/count/:user_id", follow.GetFollowCountHandler)
		authGroup.GET("/follower/count/:user_id", follow.GetFollowerCountHandler)
		authGroup.GET("/follow/mutual", follow.GetMutualFollowsHandler)
		authGroup.POST("/follow/check", follow.IsFollowingHandler)

		// 通知相关
		authGroup.GET("/notifications", handler.GetNotificationList)
	}
}