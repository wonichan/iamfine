package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	"hupu/kitex_gen/comment/commentservice"
	"hupu/services/comment/handler"
	"hupu/shared/config"
	"hupu/shared/utils"
)

func main() {
	// 初始化配置
	config.Init("../../config.yaml")

	// 初始化日志
	utils.InitLogger()

	// 初始化数据库
	db, err := utils.InitDB()
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// 初始化Redis
	rdb, err := utils.InitRedis()
	if err != nil {
		log.Fatalf("Failed to init redis: %v", err)
	}

	// 创建服务处理器
	commentHandler := handler.NewCommentHandler(db, rdb)

	// 创建服务器
	addr, _ := net.ResolveTCPAddr("tcp", config.GlobalConfig.Services.Comment.Host+":"+config.GlobalConfig.Services.Comment.Port)
	svr := commentservice.NewServer(commentHandler, server.WithServiceAddr(addr))

	utils.Logger.Info("Comment service starting...")
	err = svr.Run()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
