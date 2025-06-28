package user

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/kitex_gen/user"
)

// FollowUser 关注用户
// POST /api/user/{id}/follow
func FollowUser(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, ok := RequireAuth(c)
	if !ok {
		return
	}

	// 获取要关注的用户ID
	targetUserID, ok := ValidateTargetUserIDParam(c, "id")
	if !ok {
		return
	}

	// 不能关注自己
	if !CheckSelfFollow(c, userID, targetUserID) {
		return
	}

	// 构建请求
	req := &user.FollowUserRequest{
		UserId:       userID,
		TargetUserId: targetUserID,
	}

	// 调用用户服务
	resp, err := handler.GetUserClient().FollowUser(ctx, req)
	if err != nil {
		HandleServiceError(c, "FollowUser", err, MsgFollowUserFailed)
		return
	}

	SuccessResponse(c, resp)
}

// UnfollowUser 取消关注用户
// DELETE /api/user/{id}/follow
func UnfollowUser(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, ok := RequireAuth(c)
	if !ok {
		return
	}

	// 获取要取消关注的用户ID
	targetUserID, ok := ValidateTargetUserIDParam(c, "id")
	if !ok {
		return
	}

	// 构建请求
	req := &user.UnfollowUserRequest{
		UserId:       userID,
		TargetUserId: targetUserID,
	}

	// 调用用户服务
	resp, err := handler.GetUserClient().UnfollowUser(ctx, req)
	if err != nil {
		HandleServiceError(c, "UnfollowUser", err, MsgUnfollowUserFailed)
		return
	}

	SuccessResponse(c, resp)
}

// GetFollowers 获取粉丝列表
// GET /api/user/{id}/followers
func GetFollowers(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID, ok := ValidateUserIDParam(c, "id")
	if !ok {
		return
	}

	// 解析分页参数
	page, pageSize := ParsePaginationParams(c)

	// 构建请求
	req := &user.GetFollowersRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	// 调用用户服务
	resp, err := handler.GetUserClient().GetFollowers(ctx, req)
	if err != nil {
		HandleServiceError(c, "GetFollowers", err, MsgGetFollowersFailed)
		return
	}

	SuccessResponse(c, resp)
}

// GetFollowing 获取关注列表
// GET /api/user/{id}/following
func GetFollowing(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID, ok := ValidateUserIDParam(c, "id")
	if !ok {
		return
	}

	// 解析分页参数
	page, pageSize := ParsePaginationParams(c)

	// 构建请求
	req := &user.GetFollowingRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	// 调用用户服务
	resp, err := handler.GetUserClient().GetFollowing(ctx, req)
	if err != nil {
		HandleServiceError(c, "GetFollowing", err, MsgGetFollowingFailed)
		return
	}

	SuccessResponse(c, resp)
}
