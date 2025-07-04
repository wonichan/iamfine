package main

import (
	"context"
	"hupu/api-gateway/middleware"
	"hupu/api-gateway/router"
	"hupu/shared/config"
	"hupu/shared/log"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/logger/accesslog"
)

func main() {
	// 初始化配置
	config.Init("../../config.yaml")

	// 初始化日志
	log.InitLogger("api-gateway.log", config.GlobalConfig.Log.Path, config.GlobalConfig.Log.Level)

	// 创建Hertz服务器
	h := server.Default(
		server.WithHostPorts(config.GlobalConfig.Server.Host+":"+config.GlobalConfig.Server.Port),
		server.WithReadTimeout(60*time.Second),
		server.WithWriteTimeout(60*time.Second),
	)
	h.Use(
		accesslog.New(accesslog.WithFormat("[${time}] ${status} - ${latency} ${method} ${path} ${queryParams}")),
		middleware.TraceIdMiddleware(),
	)
	// 注册路由
	router.RegisterRoutes(h)
	h.Engine.OnRun = append(h.Engine.OnRun, beginStart)

	log.GetLogger().Info("API Gateway starting...")
	h.Spin()
}

func beginStart(ctx context.Context) error {
	log.GetLogger().Infof("here we go...")
	return nil
}
