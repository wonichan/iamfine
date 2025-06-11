package middleware

import (
	"context"
	"hupu/shared/constants"

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
