package user

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/user"
)

// CreateAnonymousProfileRequest 创建匿名马甲请求结构
type CreateAnonymousProfileRequest struct {
	AnonymousName string `json:"anonymous_name" binding:"required"`
	AvatarColor   string `json:"avatar_color"`
}

// CreateAnonymousProfile 创建匿名马甲
// POST /api/user/anonymous-profiles
func CreateAnonymousProfile(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, ok := common.RequireAuth(c)
	if !ok {
		return
	}

	// 解析请求体
	var reqBody CreateAnonymousProfileRequest
	if err := c.BindJSON(&reqBody); err != nil {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, MsgRequestFormatError)
		return
	}

	// 构建请求
	req := &user.CreateAnonymousProfileRequest{
		UserId:        userID,
		AnonymousName: reqBody.AnonymousName,
		AvatarColor:   reqBody.AvatarColor,
	}

	// 调用用户服务
	resp, err := handler.GetUserClient().CreateAnonymousProfile(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "CreateAnonymousProfile", err, MsgCreateAnonymousFailed)
		return
	}

	common.SuccessResponseFunc(c, resp)
}

// GetAnonymousProfiles 获取匿名马甲列表
// GET /api/user/anonymous-profiles
func GetAnonymousProfiles(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, ok := common.RequireAuth(c)
	if !ok {
		return
	}

	// 构建请求
	req := &user.GetAnonymousProfilesRequest{
		UserId: userID,
	}

	// 调用用户服务
	resp, err := handler.GetUserClient().GetAnonymousProfiles(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "GetAnonymousProfiles", err, MsgGetAnonymousListFailed)
		return
	}

	common.SuccessResponseFunc(c, resp)
}

// UpdateAnonymousProfile 更新匿名马甲
// PUT /api/user/anonymous-profiles/{profile_id}
func UpdateAnonymousProfile(ctx context.Context, c *app.RequestContext) {
	// 获取profile_id
	profileID, ok := common.ValidateProfileIDParam(c, "profile_id")
	if !ok {
		return
	}

	// 构建请求
	req := &user.UpdateAnonymousProfileRequest{
		ProfileId: profileID,
	}

	// 解析可选参数
	if anonymousName := c.Query("anonymous_name"); anonymousName != "" {
		req.AnonymousName = &anonymousName
	}
	if avatarColor := c.Query("avatar_color"); avatarColor != "" {
		req.AvatarColor = &avatarColor
	}
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			req.IsActive = &isActive
		}
	}

	// 调用用户服务
	resp, err := handler.GetUserClient().UpdateAnonymousProfile(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "UpdateAnonymousProfile", err, MsgUpdateAnonymousFailed)
		return
	}

	common.SuccessResponseFunc(c, resp)
}
