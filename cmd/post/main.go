package main

import (
	"context"
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"

	post "hupu/kitex_gen/post/postservice"
	"hupu/services/post/handler"
	"hupu/shared/config"
	"hupu/shared/log"
	"hupu/shared/middleware"
	"hupu/shared/utils"
)

func main() {
	// 初始化配置
	config.Init("../../config.yaml")

	// 初始化日志
	log.InitLogger("post.log", config.GlobalConfig.Log.Path, config.GlobalConfig.Log.Level)

	// 初始化数据库
	err := utils.InitDB()
	if err != nil {
		log.GetLogger().Fatalf("Failed to init database: %v", err)
		panic(err)
	}

	// 初始化Redis
	err = utils.NewRedisClient(context.Background())
	if err != nil {
		log.GetLogger().Fatalf("Failed to init redis: %v", err)
		panic(err)
	}

	// 创建服务处理器
	postHandler := handler.NewPostHandler()

	// 创建服务器
	addr, _ := net.ResolveTCPAddr("tcp", config.GlobalConfig.Services.Post.Host+":"+config.GlobalConfig.Services.Post.Port)
	svr := post.NewServer(postHandler,
		server.WithServiceAddr(addr),
		server.WithReadWriteTimeout(time.Minute),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "post"}),
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler),
		server.WithMiddleware(middleware.TraceIdMiddleWare()),
	)

	log.GetLogger().Info("Post service starting...")
	err = svr.Run()
	if err != nil {
		log.GetLogger().Fatalf("Failed to start server: %v", err)
	}
}
