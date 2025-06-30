package middleware

import (
	"context"
	"hupu/shared/constants"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/rs/xid"
)

func TraceIdMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 生成 traceId 并设置到 context
		traceId := xid.New().String()
		newCtx := metainfo.WithPersistentValue(ctx, constants.TraceIdKey, traceId)
		c.Set(constants.TraceIdKey, traceId)
		c.Next(newCtx)
	}
}
