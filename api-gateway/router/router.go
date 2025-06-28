package router

import (
	"context"

	"hupu/api-gateway/handler"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(h *server.Hertz) {
	handler.Init()

	// 健康检查
	h.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]interface{}{
			"status":  "ok",
			"message": "API Gateway is running",
		})
	})

	// 注册各模块路由
	RegisterUserRoutes(h)
	RegisterPostRoutes(h)
	RegisterSocialRoutes(h)
}
