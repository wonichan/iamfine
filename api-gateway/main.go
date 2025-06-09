package main

import (
	"hupu/api-gateway/router"
	"hupu/shared/config"
	"hupu/shared/utils"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	// 初始化配置
	config.Init("../config.yaml")

	// 初始化日志
	utils.InitLogger("api-gateway", config.GlobalConfig.Log.Path, config.GlobalConfig.Log.Level)

	// 创建Hertz服务器
	h := server.Default(server.WithHostPorts(config.GlobalConfig.Server.Host + ":" + config.GlobalConfig.Server.Port))

	// 注册路由
	router.RegisterRoutes(h)

	utils.Logger.Info("API Gateway starting...")
	h.Spin()
}
