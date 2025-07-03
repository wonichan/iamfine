package follow

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/follow"
	"hupu/shared/constants"
	"hupu/shared/log"
)

// Follow 关注用户
func Follow(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] Follow request started", traceId)

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
		common.HandleRpcError(c, "Follow", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "Follow", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// Unfollow 取消关注
func Unfollow(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] Unfollow request started", traceId)

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
		common.HandleRpcError(c, "Unfollow", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "Unfollow", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// GetFollowList 获取关注列表
func GetFollowList(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] GetFollowList request started", traceId)

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

	resp, err := handler.GetFollowClient().GetFollowList(ctx, &req)
	if err != nil {
		common.HandleRpcError(c, "GetFollowList", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetFollowList", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// GetFollowerList 获取粉丝列表
func GetFollowerList(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] GetFollowerList request started", traceId)

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

	resp, err := handler.GetFollowClient().GetFollowerList(ctx, &req)
	if err != nil {
		common.HandleRpcError(c, "GetFollowerList", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetFollowerList", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// CheckFollowStatus 检查关注状态
func CheckFollowStatus(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] CheckFollowStatus request started", traceId)

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
	resp, err := handler.GetFollowClient().CheckFollowStatus(ctx, &req)
	if err != nil {
		common.HandleRpcError(c, "CheckFollowStatus", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "CheckFollowStatus", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// IsFollowing 检查是否关注
func IsFollowing(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] IsFollowing request started", traceId)

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
		common.HandleRpcError(c, "IsFollowing", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "IsFollowing", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// GetFollowCount 获取关注数量
func GetFollowCount(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] GetFollowCount request started", traceId)

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
	resp, err := handler.GetFollowClient().GetFollowCount(ctx, &req)
	if err != nil {
		common.HandleRpcError(c, "GetFollowCount", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetFollowCount", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// GetFollowerCount 获取粉丝数量
func GetFollowerCount(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] GetFollowerCount request started", traceId)

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
	resp, err := handler.GetFollowClient().GetFollowerCount(ctx, &req)
	if err != nil {
		common.HandleRpcError(c, "GetFollowerCount", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetFollowerCount", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// GetMutualFollows 获取共同关注
func GetMutualFollows(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] GetMutualFollows request started", traceId)

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
	resp, err := handler.GetFollowClient().GetMutualFollows(ctx, &req)
	if err != nil {
		common.HandleRpcError(c, "GetMutualFollows", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetMutualFollows", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}
