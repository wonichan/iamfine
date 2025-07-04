package main

import (
	"context"
	"net"
	"time"

	"hupu/kitex_gen/comment/commentservice"
	"hupu/services/comment/handler"
	"hupu/shared/config"
	"hupu/shared/log"
	"hupu/shared/middleware"
	"hupu/shared/utils"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
)

func main() {
	// 初始化配置
	config.Init("../../config.yaml")

	// 初始化日志
	log.InitLogger("comment", config.GlobalConfig.Log.Path, config.GlobalConfig.Log.Level)

	// 初始化数据库
	if err := utils.InitDB(); err != nil {
		log.GetLogger().Fatalf("Failed to init database: %v", err)
	}

	// 初始化Redis
	if err := utils.NewRedisClient(context.Background()); err != nil {
		log.GetLogger().Fatalf("Failed to init redis: %v", err)
	}

	// 创建服务处理器
	commentHandler := handler.NewCommentHandler()

	// 创建服务器
	addr, _ := net.ResolveTCPAddr("tcp", config.GlobalConfig.Services.Comment.Host+":"+config.GlobalConfig.Services.Comment.Port)
	svr := commentservice.NewServer(commentHandler,
		server.WithServiceAddr(addr),
		server.WithReadWriteTimeout(time.Minute),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "comment"}),
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler),
		server.WithMiddleware(middleware.TraceIdMiddleWare()),
	)

	log.GetLogger().Info("Comment service starting...")
	err := svr.Run()
	if err != nil {
		log.GetLogger().Fatalf("Failed to start server: %v", err)
	}
}
