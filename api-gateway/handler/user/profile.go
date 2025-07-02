package user

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/user"
	"hupu/shared/constants"
	"hupu/shared/log"
	"hupu/shared/models"
)

// GetUserStats 获取用户统计
// GET /api/user/stats
func GetUserStats(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	// 需要认证
	userID, ok := common.RequireAuth(c)
	if !ok {
		common.RespondUnauthorized(c)
		return
	}

	// 调用用户服务获取统计信息
	req := &user.GetUserStatsRequest{
		UserId: userID,
	}

	resp, err := handler.GetUserClient().GetUserStats(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "GetUserStats", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetUserStats", traceId.(string), resp.Code, resp.Message)
		return
	}

	// 转换响应格式以符合API文档
	responseData := map[string]interface{}{
		"code":    common.CodeSuccess,
		"message": constants.MsgSuccess,
		"data": map[string]interface{}{
			"posts":     resp.PostCount,
			"comments":  resp.CommentCount,
			"favorites": resp.LikeCount,
			"ratings":   resp.AverageScore,
			"followers": resp.FollowerCount,
			"following": resp.FollowingCount,
		},
	}

	common.SuccessResponseFunc(c, responseData)
}

// GetUnreadCount 获取未读消息数
// GET /api/user/unread-count
func GetUnreadCount(ctx context.Context, c *app.RequestContext) {
	// 需要认证
	_, ok := common.RequireAuth(c)
	if !ok {
		common.RespondUnauthorized(c)
		return
	}

	// TODO: 调用通知服务获取未读消息数
	// 这里需要调用notification服务
	// 由于当前没有相应的thrift接口，先返回示例数据

	responseData := map[string]interface{}{
		"code":    common.CodeSuccess,
		"message": constants.MsgSuccess,
		"data": map[string]interface{}{
			"count": 0,
		},
	}

	common.SuccessResponseFunc(c, responseData)
}

// GetUser 获取指定用户信息（保留原有接口）
// GET /api/user/{id}
func GetUser(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	// 获取用户ID参数
	userID, ok := common.ValidateUserIDParam(c, "id")
	if !ok {
		logger.Errorf("GetUser ValidateUserIDParam failed")
		common.RespondBadRequest(c, constants.MsgRequestFormatError)
		return
	}

	// 调用用户服务
	req := &user.GetUserRequest{
		UserId: userID,
	}
	resp, err := handler.GetUserClient().GetUser(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "GetUser", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetUser", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// UpdateUser 更新指定用户信息（保留原有接口）
// PUT /api/user/{id}
func UpdateUser(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	// 从上下文获取用户ID
	currentUserID, ok := common.RequireAuth(c)
	if !ok {
		logger.Errorf("UpdateUser RequireAuth failed")
		common.RespondUnauthorized(c)
		return
	}

	// 获取要更新的用户ID
	userID, ok := common.ValidateUserIDParam(c, "id")
	if !ok {
		logger.Errorf("UpdateUser ValidateUserIDParam failed")
		common.RespondBadRequest(c, constants.MsgRequestFormatError)
		return
	}

	// 检查权限：只能更新自己的信息
	if currentUserID != userID {
		logger.Errorf("UpdateUser no permission")
		common.ErrorResponseFunc(c, constants.HTTPStatusForbidden, common.CodeError, constants.MsgNoPermissionUpdateOther)
		return
	}

	// 解析请求体
	var reqBody models.User

	if err := c.BindJSON(&reqBody); err != nil {
		logger.Errorf("UpdateUser BindJSON failed, err:%s", err.Error())
		common.RespondBadRequest(c, constants.MsgRequestFormatError)
		return
	}

	// 构建请求
	req := &user.UpdateUserRequest{
		Id: userID,
	}

	if reqBody.Nickname != "" {
		req.Nickname = &reqBody.Nickname
	}
	if reqBody.Avatar != "" {
		req.Avatar = &reqBody.Avatar
	}
	if reqBody.Bio != nil {
		req.Bio = reqBody.Bio
	}
	if reqBody.Location != nil {
		req.Location = reqBody.Location
	}
	if reqBody.RelationshipStatus != nil {
		req.RelationshipStatus = reqBody.RelationshipStatus
	}
	if reqBody.Tags != nil {
		req.Tags = reqBody.Tags
	}
	if reqBody.Email != "" {
		req.Email = &reqBody.Email
	}
	if reqBody.Phone != "" {
		req.Phone = &reqBody.Phone
	}

	// 调用用户服务
	resp, err := handler.GetUserClient().UpdateUser(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "UpdateUser", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "UpdateUser", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}
