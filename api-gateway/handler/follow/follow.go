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
	resp, err := handler.GetFollowClient().Follow(ctx, &req)
	if err != nil {
		common.HandleServiceError(c, "Follow", err, constants.MsgFollowFailed)
		return
	}
	common.RespondWithSuccess(c, resp)
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
	resp, err := handler.GetFollowClient().Unfollow(ctx, &req)
	if err != nil {
		common.HandleServiceError(c, "Unfollow", err, constants.MsgUnfollowFailed)
		return
	}

	common.RespondWithSuccess(c, resp)
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
	req := &follow.GetFollowListRequest{
		UserId:   userID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}
	resp, err := handler.GetFollowClient().GetFollowList(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "GetFollowList", err, constants.MsgGetFollowListFailed)
		return
	}

	common.RespondWithSuccess(c, resp)
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
	req := &follow.GetFollowerListRequest{
		UserId:   userID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}
	resp, err := handler.GetFollowClient().GetFollowerList(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "GetFollowerList", err, constants.MsgGetFollowerListFailed)
		return
	}

	common.RespondWithSuccess(c, resp)
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
	req := &follow.CheckFollowStatusRequest{
		FollowerId:  userID,
		FollowingId: targetUserID,
	}

	// 调用关注服务
	resp, err := handler.GetFollowClient().CheckFollowStatus(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "GetFollowerList", err, constants.MsgCheckFollowStatusFailed)
		return
	}

	common.RespondWithSuccess(c, resp)
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
	resp, err := handler.GetFollowClient().IsFollowing(ctx, &req)
	if err != nil {
		common.HandleServiceError(c, "IsFollowing", err, constants.MsgCheckFollowStatusFailed)
		return
	}

	common.RespondWithSuccess(c, resp)
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
	req := &follow.GetFollowCountRequest{
		UserId: userID,
	}

	// 调用关注服务
	resp, err := handler.GetFollowClient().GetFollowCount(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "GetFollowCount", err, constants.MsgGetFollowCountFailed)
		return
	}

	common.RespondWithSuccess(c, resp)
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
	req := &follow.GetFollowerCountRequest{
		UserId: userID,
	}

	// 调用关注服务
	resp, err := handler.GetFollowClient().GetFollowerCount(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "GetFollowerCount", err, constants.MsgGetFollowerCountFailed)
		return
	}

	common.RespondWithSuccess(c, resp)
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
	req := &follow.GetMutualFollowsRequest{
		UserId:       userID,
		TargetUserId: targetUserID,
		Page:         int32(page),
		PageSize:     int32(pageSize),
	}

	// 调用关注服务
	resp, err := handler.GetFollowClient().GetMutualFollows(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "GetMutualFollows", err, constants.MsgGetMutualFollowsFailed)
		return
	}

	common.RespondWithSuccess(c, resp)
}
