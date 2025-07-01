package main

import (
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"

	user "hupu/kitex_gen/user/userservice"
	"hupu/services/user/handler"
	"hupu/shared/config"
	"hupu/shared/log"
	"hupu/shared/middleware"
	"hupu/shared/utils"
)

func main() {
	// 初始化配置
	config.Init("../../config.yaml")

	// 初始化日志
	log.InitLogger("user", config.GlobalConfig.Log.Path, config.GlobalConfig.Log.Level)

	// 初始化数据库
	db, err := utils.InitDB()
	if err != nil {
		log.GetLogger().Errorf("Failed to init database: %v", err)
		panic("Failed to init database")
	}

	// 初始化Redis
	// rdb, err := utils.InitRedis()
	// if err != nil {
	// 	log.GetLogger().Fatalf("Failed to init redis: %v", err)
	// }

	// 创建服务处理器
	userHandler := handler.NewUserHandler(db)

	// 创建服务器
	addr, _ := net.ResolveTCPAddr("tcp", config.GlobalConfig.Services.User.Host+":"+config.GlobalConfig.Services.User.Port)
	svr := user.NewServer(userHandler,
		server.WithServiceAddr(addr),
		server.WithReadWriteTimeout(time.Minute),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "user"}),
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler),
		server.WithMiddleware(middleware.TraceIdMiddleWare()),
	)

	log.GetLogger().Info("User service starting...")
	err = svr.Run()
	if err != nil {
		log.GetLogger().Fatalf("Failed to start server: %v", err)
	}
}
