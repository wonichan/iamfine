package user

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/user"
	"hupu/shared/constants"
)

// FollowUser 关注用户
// POST /api/user/{id}/follow
func FollowUser(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, ok := common.RequireAuth(c)
	if !ok {
		return
	}

	// 获取要关注的用户ID
	targetUserID, ok := common.ValidateTargetUserIDParam(c, "id")
	if !ok {
		return
	}

	// 不能关注自己
	if userID == targetUserID {
		common.RespondBadRequest(c, "不能关注自己")
		return
	}
	if false {
		return
	}

	// 构建请求
	req := &user.FollowUserRequest{
		UserId:       userID,
		TargetUserId: targetUserID,
	}

	// 调用用户服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetUserClient().FollowUser(ctx, req)
	}), "FollowUser", constants.MsgFollowUserFailed)
}

// UnfollowUser 取消关注用户
// DELETE /api/user/{id}/follow
func UnfollowUser(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, ok := common.RequireAuth(c)
	if !ok {
		return
	}

	// 获取要取消关注的用户ID
	targetUserID, ok := common.ValidateTargetUserIDParam(c, "id")
	if !ok {
		return
	}

	// 构建请求
	req := &user.UnfollowUserRequest{
		UserId:       userID,
		TargetUserId: targetUserID,
	}

	// 调用用户服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetUserClient().UnfollowUser(ctx, req)
	}), "UnfollowUser", constants.MsgUnfollowUserFailed)
}

// GetFollowers 获取粉丝列表
// GET /api/user/{id}/followers
func GetFollowers(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID, ok := common.ValidateUserIDParam(c, "id")
	if !ok {
		return
	}

	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 构建请求
	req := &user.GetFollowersRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	// 调用用户服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetUserClient().GetFollowers(ctx, req)
	}), "GetFollowers", constants.MsgGetFollowersFailed)
}

// GetFollowing 获取关注列表
// GET /api/user/{id}/following
func GetFollowing(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID, ok := common.ValidateUserIDParam(c, "id")
	if !ok {
		return
	}

	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 构建请求
	req := &user.GetFollowingRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	// 调用用户服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetUserClient().GetFollowing(ctx, req)
	}), "GetFollowing", constants.MsgGetFollowingFailed)
}
