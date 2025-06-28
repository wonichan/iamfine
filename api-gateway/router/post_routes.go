package router

import (
	"hupu/api-gateway/handler/post"
	"hupu/api-gateway/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterPostRoutes 注册帖子相关路由
func RegisterPostRoutes(h *server.Hertz) {
	// API版本组
	apiV1 := h.Group("/api/v1")

	// 帖子相关路由
	postGroup := apiV1.Group("/posts")
	{
		// 无需认证的路由
		postGroup.GET("/", post.GetPostListHandler)
		postGroup.GET("/:id", post.GetPostHandler)
		postGroup.GET("/recommend", post.GetRecommendPostsHandler)
		postGroup.GET("/hot", post.GetHotPostsHandler)
		postGroup.GET("/high-score", post.GetHighScorePostsHandler)
		postGroup.GET("/low-score", post.GetLowScorePostsHandler)
		postGroup.GET("/controversial", post.GetControversialPostsHandler)
		postGroup.GET("/search", post.SearchPostsHandler)

		// 需要认证的路由
		postAuthGroup := postGroup.Group("/", middleware.AuthMiddleware())
		{
			postAuthGroup.POST("/", post.CreatePostHandler)
			postAuthGroup.PUT("/:id", post.UpdatePostHandler)
			postAuthGroup.DELETE("/:id", post.DeletePostHandler)
			postAuthGroup.POST("/:id/collect", post.CollectPostHandler)
			postAuthGroup.DELETE("/:id/collect", post.UncollectPostHandler)
			postAuthGroup.GET("/collected", post.GetCollectedPostsHandler)
			postAuthGroup.POST("/:id/rate", post.RatePostHandler)
			postAuthGroup.GET("/:id/rating", post.GetUserRatingHandler)
			postAuthGroup.PUT("/:id/rating", post.UpdateRatingHandler)
			postAuthGroup.DELETE("/:id/rating", post.DeleteRatingHandler)
			postAuthGroup.GET("/rating/rank", post.GetRatingRankHandler)
		}

		// 话题相关路由
		topicGroup := postGroup.Group("/topics")
		{
			// 无需认证的路由
			topicGroup.GET("/", post.GetTopicListHandler)
			topicGroup.GET("/:id", post.GetTopicHandler)
			topicGroup.GET("/hot", post.GetHotTopicsHandler)
			topicGroup.GET("/categories", post.GetTopicCategoriesHandler)
			topicGroup.GET("/search", post.SearchTopicsHandler)

			// 需要认证的路由
			topicAuthGroup := topicGroup.Group("/", middleware.AuthMiddleware())
			{
				topicAuthGroup.POST("/", post.CreateTopicHandler)
				topicAuthGroup.POST("/:id/share", post.ShareTopicHandler)
			}
		}
	}
}