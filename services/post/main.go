package main

import (
	"context"
	"github.com/cloudwego/kitex/server"
	post "hupu/kitex_gen/post/postservice"
	"hupu/services/post/handler"
	"hupu/shared/config"
	"hupu/shared/utils"
	"log"
	"net"
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
	postHandler := handler.NewPostHandler(db, rdb)

	// 创建服务器
	addr, _ := net.ResolveTCPAddr("tcp", config.GlobalConfig.Services.Post.Host+":"+config.GlobalConfig.Services.Post.Port)
	svr := post.NewServer(postHandler, server.WithServiceAddr(addr))

	utils.Logger.Info("Post service starting...")
	err = svr.Run()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
