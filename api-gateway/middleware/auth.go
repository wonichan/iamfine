package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/rs/xid"

	"hupu/shared/constants"
	"hupu/shared/log"
	"hupu/shared/utils"
)

func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 获取Authorization header
		auth := c.GetHeader("Authorization")
		if len(auth) == 0 {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "未提供认证token",
			})
			c.Abort()
			return
		}

		// 解析Bearer token
		tokenString := string(auth)
		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "token格式错误",
			})
			c.Abort()
			return
		}

		tokenString = tokenString[7:] // 移除"Bearer "

		// 验证token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			log.GetLogger().Errorf("Parse token error: %v", err)
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "token无效",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next(ctx)
	}
}

func TraceIdMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成 traceId 并设置到 context
		traceId := xid.New().String()
		newCtx := metainfo.WithPersistentValue(ctx, constants.TraceIdKey, traceId)
		c.Next(newCtx)
	}
}
