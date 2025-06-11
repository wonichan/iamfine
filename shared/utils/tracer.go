package utils

import (
	"context"
	"fmt"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/stats"

	"hupu/shared/constants"
	"hupu/shared/log"
)

// KitexClientTracer 客户端链路追踪器
type KitexClientTracer struct {
}

// NewKitexClientTracer 创建一个新的客户端链路追踪器
func NewKitexClientTracer() *KitexClientTracer {
	return &KitexClientTracer{}
}

func (t *KitexClientTracer) Start(ctx context.Context) context.Context {
	return ctx
}

func (t *KitexClientTracer) Finish(ctx context.Context) {
	// 获取traceID
	traceID, ok := metainfo.GetPersistentValue(ctx, constants.TraceIdKey)
	if traceID == "" || !ok {
		klog.Warnf("trace_id not found in context")
		return
	}

	// 获取RPC信息
	ri := rpcinfo.GetRPCInfo(ctx)
	if ri == nil {
		klog.Warnf("rpcinfo not found in context")
		return
	}

	rpcStart := ri.Stats().GetEvent(stats.RPCStart)
	rpcFinish := ri.Stats().GetEvent(stats.RPCFinish)
	cost := rpcFinish.Time().Sub(rpcStart.Time())

	// 获取请求和响应的大小
	reqSize := ri.Stats().RecvSize()
	respSize := ri.Stats().SendSize()

	// 获取服务和方法信息
	service := ri.To().ServiceName()
	method := ri.To().Method()

	// 记录请求信息
	logMsg := fmt.Sprintf(
		"TraceID: %s, Service: %s, Method: %s, ReqSize: %d bytes, RespSize: %d bytes, Time: %v",
		traceID, service, method, reqSize, respSize, cost,
	)

	// 记录请求信息
	log.GetLogger().Infof(logMsg)
}
