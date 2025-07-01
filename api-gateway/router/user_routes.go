package router

import (
	"hupu/api-gateway/handler/user"
	"hupu/api-gateway/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(h *server.Hertz) {
	// API版本组
	apiV1 := h.Group("/api/v1")

	// 用户相关路由（无需认证）
	userGroup := apiV1.Group("/user")
	{
		userGroup.POST("/register", user.RegisterHandler)
		userGroup.POST("/login", user.LoginHandler)
		userGroup.POST("/wx-login", user.WxLoginHandler)
		userGroup.GET("/profile/:id", user.GetUserHandler)
	}

	// 需要认证的用户路由
	authGroup := apiV1.Group("", middleware.AuthMiddleware())
	{
		// 用户相关（需要认证）
		authGroup.GET("/user/info", user.GetUserInfoHandler)
		authGroup.GET("/user/stats", user.GetUserStatsHandler)
		authGroup.GET("/user/unread-count", user.GetUnreadCountHandler)
		authGroup.POST("/user/:id/follow", user.FollowUserHandler)
		authGroup.DELETE("/user/:id/follow", user.UnfollowUserHandler)
		authGroup.GET("/user/:id/followers", user.GetFollowersHandler)
		authGroup.GET("/user/:id/following", user.GetFollowingHandler)
		authGroup.POST("/user/anonymous-profiles", user.CreateAnonymousProfileHandler)
		authGroup.GET("/user/anonymous-profiles", user.GetAnonymousProfilesHandler)
		authGroup.PUT("/user/anonymous-profiles/:profile_id", user.UpdateAnonymousProfileHandler)
		authGroup.PUT("/user/:id", user.UpdateUserHandler)
	}
}
