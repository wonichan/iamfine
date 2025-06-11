package main

import (
	"net"
	"time"

	"hupu/kitex_gen/like/likeservice"
	"hupu/services/like/handler"
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
	log.InitLogger("like", config.GlobalConfig.Log.Path, config.GlobalConfig.Log.Level)

	// 初始化数据库
	db, err := utils.InitDB()
	if err != nil {
		log.GetLogger().Fatalf("Failed to init database: %v", err)
	}

	// 初始化Redis
	rdb, err := utils.InitRedis()
	if err != nil {
		log.GetLogger().Fatalf("Failed to init redis: %v", err)
	}

	// 创建服务处理器
	likeHandler := handler.NewLikeHandler(db, rdb)

	// 创建服务器
	addr, _ := net.ResolveTCPAddr("tcp", config.GlobalConfig.Services.Like.Host+":"+config.GlobalConfig.Services.Like.Port)
	svr := likeservice.NewServer(likeHandler,
		server.WithServiceAddr(addr),
		server.WithReadWriteTimeout(time.Minute),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "like"}),
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler),
		server.WithMiddleware(middleware.TraceIdMiddleWare()),
	)

	log.GetLogger().Info("Like service starting...")
	err = svr.Run()
	if err != nil {
		log.GetLogger().Fatalf("Failed to start server: %v", err)
	}
}
