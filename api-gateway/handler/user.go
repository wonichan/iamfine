package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/kitex_gen/user"
	"hupu/shared/utils"
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

	// 调用用户服务
	resp, err := userClient.Register(ctx, &req)
	if err != nil {
		utils.GetLogger().Errorf("Register error: %v", err)
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

	// 调用用户服务
	resp, err := userClient.Login(ctx, &req)
	if err != nil {
		utils.GetLogger().Errorf("Login error: %v", err)
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

	// 调用用户服务
	req := &user.GetUserRequest{UserId: userID}
	resp, err := userClient.GetUser(ctx, req)
	if err != nil {
		utils.GetLogger().Errorf("GetUser error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
