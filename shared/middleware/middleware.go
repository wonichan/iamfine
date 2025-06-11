package middleware

import (
	"context"
	"hupu/shared/constants"

	"github.com/rs/xid"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/pkg/endpoint"
)

func TraceIdMiddleWare() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			traceId, _ := metainfo.GetPersistentValue(ctx, constants.TraceIdKey)

			ctx = context.WithValue(ctx, constants.TraceIdKey, traceId)
			return next(ctx, req, resp)
		}
	}
}

func WithTraceIdMiddleWare() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) (err error) {
			traceId := ctx.Value(constants.TraceIdKey)
			if traceId == nil {
				traceId = xid.New().String()
			}
			ctx = metainfo.WithPersistentValue(ctx, constants.TraceIdKey, traceId.(string))
			return next(ctx, req, resp)
		}
	}
}
