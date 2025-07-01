package user

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/user"
	"hupu/shared/constants"
	"hupu/shared/log"
)

// WxLoginRequest 微信登录请求结构
type WxLoginRequest struct {
	Code string `json:"code" binding:"required"` // 微信登录凭证
}

// WxLogin 微信登录
// POST /api/user/wx-login
func WxLogin(ctx context.Context, c *app.RequestContext) {
	// 解析请求参数
	var req WxLoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		common.RespondBadRequest(c, constants.MsgParamError+": "+err.Error())
		return
	}

	// TODO: 调用微信登录服务
	// 这里需要根据实际的微信登录服务接口进行调用
	// 由于当前thrift定义中没有微信登录接口，这里先返回示例响应

	// 示例响应数据
	responseData := map[string]interface{}{
		"code":    common.CodeSuccess,
		"message": constants.MsgSuccess,
		"data": map[string]interface{}{
			"token": "example_jwt_token",
			"user": map[string]interface{}{
				"id":        "user_id_example",
				"openId":    "wx_open_id_example",
				"unionId":   "wx_union_id_example",
				"avatar":    "https://example.com/avatar.jpg",
				"nickname":  "微信用户",
				"gender":    0,
				"city":      "深圳",
				"province":  "广东",
				"country":   "中国",
				"tags":      []string{},
				"createdAt": "2024-01-01T00:00:00Z",
				"updatedAt": "2024-01-01T00:00:00Z",
			},
		},
	}

	common.SuccessResponseFunc(c, responseData)
}

// GetUserInfo 获取用户信息
// GET /api/user/info
func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	// 需要认证
	userID, ok := common.RequireAuth(c)
	if !ok {
		logger.Errorf("GetUserInfo RequireAuth failed")
		common.RespondUnauthorized(c)
		return
	}

	// 调用用户服务获取用户信息
	req := &user.GetUserRequest{
		UserId: userID,
	}

	resp, err := handler.GetUserClient().GetUser(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "GetUserInfo", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetUserInfo", traceId.(string), resp.Code, resp.Message)
		return
	}

	common.SuccessResponseFunc(c, resp.User)
}
