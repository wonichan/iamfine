package follow

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/follow"
	"hupu/shared/constants"
)

// Follow 关注用户
func Follow(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 解析请求参数
	var req follow.FollowRequest
	if err := c.BindAndValidate(&req); err != nil {
		common.RespondBadRequest(c, constants.MsgParamError+": "+err.Error())
		return
	}

	// 设置关注者ID
	req.FollowerId = userID

	// 调用关注服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetFollowClient().Follow(ctx, &req)
	}), "Follow", constants.MsgFollowFailed)
}

// Unfollow 取消关注
func Unfollow(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 解析请求参数
	var req follow.UnfollowRequest
	if err := c.BindAndValidate(&req); err != nil {
		common.RespondBadRequest(c, constants.MsgParamError+": "+err.Error())
		return
	}

	// 设置关注者ID
	req.FollowerId = userID

	// 调用取消关注服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetFollowClient().Unfollow(ctx, &req)
	}), "Unfollow", constants.MsgUnfollowFailed)
}

// GetFollowList 获取关注列表
func GetFollowList(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID := common.GetPathParam(c, common.UserIDKey)
	if userID == "" {
		common.RespondBadRequest(c, constants.MsgUserIDRequired)
		return
	}

	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 调用关注服务
	req := follow.GetFollowListRequest{
		UserId:   userID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetFollowClient().GetFollowList(ctx, &req)
	}), "GetFollowList", constants.MsgGetFollowListFailed)
}

// GetFollowerList 获取粉丝列表
func GetFollowerList(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID := common.GetPathParam(c, common.UserIDKey)
	if userID == "" {
		common.RespondBadRequest(c, constants.MsgUserIDRequired)
		return
	}

	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 调用关注服务
	req := follow.GetFollowerListRequest{
		UserId:   userID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetFollowClient().GetFollowerList(ctx, &req)
	}), "GetFollowerList", constants.MsgGetFollowerListFailed)
}

// CheckFollowStatus 检查关注状态
func CheckFollowStatus(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 获取目标用户ID参数
	targetUserID := c.Query("target_user_id")
	if targetUserID == "" {
		common.RespondBadRequest(c, constants.MsgTargetUserIDRequired)
		return
	}

	// 构建请求
	req := follow.CheckFollowStatusRequest{
		FollowerId:  userID,
		FollowingId: targetUserID,
	}

	// 调用关注服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetFollowClient().CheckFollowStatus(ctx, &req)
	}), "CheckFollowStatus", constants.MsgCheckFollowStatusFailed)
}

// IsFollowing 检查是否关注
func IsFollowing(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 解析请求参数
	var req follow.FollowRequest
	if err := c.BindAndValidate(&req); err != nil {
		common.RespondBadRequest(c, constants.MsgParamError+": "+err.Error())
		return
	}

	// 设置关注者ID
	req.FollowerId = userID

	// 调用关注服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetFollowClient().IsFollowing(ctx, &req)
	}), "IsFollowing", constants.MsgCheckFollowStatusFailed)
}

// GetFollowCount 获取关注数量
func GetFollowCount(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID := common.GetPathParam(c, common.UserIDKey)
	if userID == "" {
		common.RespondBadRequest(c, constants.MsgUserIDRequired)
		return
	}

	// 构建请求
	req := follow.GetFollowCountRequest{
		UserId: userID,
	}

	// 调用关注服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetFollowClient().GetFollowCount(ctx, &req)
	}), "GetFollowCount", constants.MsgGetFollowCountFailed)
}

// GetFollowerCount 获取粉丝数量
func GetFollowerCount(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID := common.GetPathParam(c, common.UserIDKey)
	if userID == "" {
		common.RespondBadRequest(c, constants.MsgUserIDRequired)
		return
	}

	// 构建请求
	req := follow.GetFollowerCountRequest{
		UserId: userID,
	}

	// 调用关注服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetFollowClient().GetFollowerCount(ctx, &req)
	}), "GetFollowerCount", constants.MsgGetFollowerCountFailed)
}

// GetMutualFollows 获取共同关注
func GetMutualFollows(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 获取目标用户ID参数
	targetUserID := c.Query("target_user_id")
	if targetUserID == "" {
		common.RespondBadRequest(c, constants.MsgTargetUserIDRequired)
		return
	}

	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 构建请求
	req := follow.GetMutualFollowsRequest{
		UserId:       userID,
		TargetUserId: targetUserID,
		Page:         int32(page),
		PageSize:     int32(pageSize),
	}

	// 调用关注服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return handler.GetFollowClient().GetMutualFollows(ctx, &req)
	}), "GetMutualFollows", constants.MsgGetMutualFollowsFailed)
}
