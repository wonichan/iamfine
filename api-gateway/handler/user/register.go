package user

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/user"
)

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname"`
}

// Register 用户注册（保留原有接口）
// POST /api/user/register
func Register(ctx context.Context, c *app.RequestContext) {
	// 解析请求体
	var reqBody RegisterRequest
	if err := c.BindJSON(&reqBody); err != nil {
		common.ErrorResponseFunc(c, common.HTTPStatusBadRequest, common.CodeError, common.MsgRequestFormatError)
		return
	}

	// 构建请求
	req := &user.RegisterRequest{
		Username: reqBody.Username,
		Password: reqBody.Password,
	}

	// 调用用户服务
	resp, err := handler.GetUserClient().Register(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "Register", err, MsgRegisterFailed)
		return
	}

	common.SuccessResponseFunc(c, resp)
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录（保留原有接口）
// POST /api/user/login
func Login(ctx context.Context, c *app.RequestContext) {
	// 解析请求体
	var reqBody LoginRequest
	if err := c.BindJSON(&reqBody); err != nil {
		common.ErrorResponseFunc(c, common.HTTPStatusBadRequest, common.CodeError, common.MsgRequestFormatError)
		return
	}

	// 构建请求
	req := &user.LoginRequest{
		Username: reqBody.Username,
		Password: reqBody.Password,
	}

	// 调用用户服务
	resp, err := handler.GetUserClient().Login(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "Login", err, MsgLoginFailed)
		return
	}

	c.Header("Authorization", resp.Token)
	common.SuccessResponseFunc(c, resp)
}
