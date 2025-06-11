package utils

import (
	"context"
	"fmt"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/rs/xid"

	"hupu/shared/constants"
)

func ContextWithTraceId(ctx context.Context, prefix string) context.Context {
	traceId := xid.New().String()
	return metainfo.WithPersistentValue(ctx, constants.TraceIdKey, fmt.Sprintf("%s-%s", prefix, traceId))
}
