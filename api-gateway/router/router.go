package router

import (
	"context"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz) {
	handler.Init()
	// 健康检查
	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]interface{}{
			"status":  "ok",
			"message": "API Gateway is running",
		})
	})

	// API版本组
	apiV1 := h.Group("/api/v1")

	// 用户相关路由（无需认证）
	userGroup := apiV1.Group("/user")
	{
		userGroup.POST("/register", handler.Register)
		userGroup.POST("/login", handler.Login)
		userGroup.GET("/profile/:id", handler.GetUser)
	}

	// 需要认证的路由
	authGroup := apiV1.Group("", middleware.AuthMiddleware())
	{
		// 帖子相关
		authGroup.POST("/post", handler.CreatePost)
		authGroup.GET("/post/:id", handler.GetPost)
		authGroup.GET("/posts", handler.GetPostList)

		// 评论相关
		authGroup.POST("/comment", handler.CreateComment)
		authGroup.GET("/comments/:post_id", handler.GetCommentList)

		// 点赞相关
		authGroup.POST("/like", handler.Like)
		authGroup.DELETE("/like", handler.Unlike)
		authGroup.GET("/likes", handler.GetLikeList) // 添加这行

		// 关注相关
		authGroup.POST("/follow", handler.Follow)
		authGroup.DELETE("/follow", handler.Unfollow)
		authGroup.GET("/follows/:user_id", handler.GetFollowList)
		authGroup.GET("/followers/:user_id", handler.GetFollowerList) // 新增获取粉丝列表

		// 通知相关
		authGroup.GET("/notifications", handler.GetNotificationList)
	}
}
