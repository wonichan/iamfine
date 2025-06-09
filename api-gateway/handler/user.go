package handler

import (
	"context"
	"net/http"

	"hupu/kitex_gen/user"
	"hupu/kitex_gen/user/userservice"
	"hupu/shared/config"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
)

// 用户注册
func Register(ctx context.Context, c *app.RequestContext) {
	// 解析请求参数
	var req user.RegisterRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 创建用户服务客户端
	client, err := userservice.NewClient("user", client.WithHostPorts(config.GlobalConfig.Services.User.Host+":"+config.GlobalConfig.Services.User.Port))
	if err != nil {
		hlog.Errorf("Create user client error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "服务连接失败",
		})
		return
	}

	// 调用用户服务
	resp, err := client.Register(ctx, &req)
	if err != nil {
		hlog.Errorf("Register error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "注册失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 用户登录
func Login(ctx context.Context, c *app.RequestContext) {
	// 解析请求参数
	var req user.LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 创建用户服务客户端
	client, err := userservice.NewClient("user", client.WithHostPorts(config.GlobalConfig.Services.User.Host+":"+config.GlobalConfig.Services.User.Port))
	if err != nil {
		hlog.Errorf("Create user client error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "服务连接失败",
		})
		return
	}

	// 调用用户服务
	resp, err := client.Login(ctx, &req)
	if err != nil {
		hlog.Errorf("Login error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "登录失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取用户信息
func GetUser(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 创建用户服务客户端
	client, err := userservice.NewClient("user", client.WithHostPorts(config.GlobalConfig.Services.User.Host+":"+config.GlobalConfig.Services.User.Port))
	if err != nil {
		hlog.Errorf("Create user client error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "服务连接失败",
		})
		return
	}

	// 调用用户服务
	req := &user.GetUserRequest{UserId: userID}
	resp, err := client.GetUser(ctx, req)
	if err != nil {
		hlog.Errorf("GetUser error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
