package user

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/user"
)

// UpdateUserInfoRequest 更新用户信息请求结构
type UpdateUserInfoRequest struct {
	Avatar   *string   `json:"avatar"`
	Nickname *string   `json:"nickname"`
	Gender   *int32    `json:"gender"`
	City     *string   `json:"city"`
	Province *string   `json:"province"`
	Country  *string   `json:"country"`
	Tags     *[]string `json:"tags"`
}

// UpdateUserInfo 更新用户信息
// PUT /api/user/info
func UpdateUserInfo(ctx context.Context, c *app.RequestContext) {
	// 需要认证
	currentUserID, ok := common.RequireAuth(c)
	if !ok {
		return
	}

	// 解析请求体
	var reqBody UpdateUserInfoRequest
	if err := c.BindJSON(&reqBody); err != nil {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, common.CodeError, MsgRequestFormatError)
		return
	}

	// 构建更新请求
	req := &user.UpdateUserRequest{
		Id: currentUserID,
	}

	if reqBody.Nickname != nil {
		req.Nickname = reqBody.Nickname
	}
	if reqBody.Avatar != nil {
		req.Avatar = reqBody.Avatar
	}
	// 注意：这里需要根据实际的thrift定义来映射字段
	// 当前thrift定义中没有gender, city, province, country字段
	// 如果需要支持这些字段，需要更新thrift定义

	if reqBody.Tags != nil {
		req.Tags = *reqBody.Tags
	}

	// 调用用户服务
	resp, err := handler.GetUserClient().UpdateUser(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "UpdateUserInfo", err, MsgUpdateUserInfoFailed)
		return
	}

	common.SuccessResponseFunc(c, resp)
}

// GetUserStats 获取用户统计
// GET /api/user/stats
func GetUserStats(ctx context.Context, c *app.RequestContext) {
	// 需要认证
	userID, ok := common.RequireAuth(c)
	if !ok {
		return
	}

	// 调用用户服务获取统计信息
	req := &user.GetUserStatsRequest{
		UserId: userID,
	}

	resp, err := handler.GetUserClient().GetUserStats(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "GetUserStats", err, MsgGetUserStatsFailed)
		return
	}

	// 转换响应格式以符合API文档
	responseData := map[string]interface{}{
		"code":    CodeSuccess,
		"message": MsgSuccess,
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
		return
	}

	// TODO: 调用通知服务获取未读消息数
	// 这里需要调用notification服务
	// 由于当前没有相应的thrift接口，先返回示例数据

	responseData := map[string]interface{}{
		"code":    CodeSuccess,
		"message": MsgSuccess,
		"data": map[string]interface{}{
			"count": 0,
		},
	}

	common.SuccessResponseFunc(c, responseData)
}

// GetUser 获取指定用户信息（保留原有接口）
// GET /api/user/{id}
func GetUser(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID, ok := common.ValidateUserIDParam(c, "id")
	if !ok {
		return
	}

	// 调用用户服务
	req := &user.GetUserRequest{
		UserId: userID,
	}
	resp, err := handler.GetUserClient().GetUser(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "GetUser", err, MsgGetUserInfoFailed)
		return
	}

	common.SuccessResponseFunc(c, resp)
}

// UpdateUser 更新指定用户信息（保留原有接口）
// PUT /api/user/{id}
func UpdateUser(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	currentUserID, ok := common.RequireAuth(c)
	if !ok {
		return
	}

	// 获取要更新的用户ID
	userID, ok := common.ValidateUserIDParam(c, "id")
	if !ok {
		return
	}

	// 检查权限：只能更新自己的信息
	if currentUserID != userID {
		common.ErrorResponseFunc(c, common.HTTPStatusForbidden, common.CodeError, "无权限修改其他用户信息")
		return
	}

	// 解析请求体
	var reqBody struct {
		Nickname  *string `json:"nickname"`
		Avatar    *string `json:"avatar"`
		Bio       *string `json:"bio"`
		Gender    *int32  `json:"gender"`
		Birthdate *string `json:"birthdate"`
		Location  *string `json:"location"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		common.ErrorResponseFunc(c, common.HTTPStatusBadRequest, common.CodeError, common.MsgRequestFormatError)
		return
	}

	// 构建请求
	req := &user.UpdateUserRequest{
		Id: userID,
	}

	if reqBody.Nickname != nil {
		req.Nickname = reqBody.Nickname
	}
	if reqBody.Avatar != nil {
		req.Avatar = reqBody.Avatar
	}
	if reqBody.Bio != nil {
		req.Bio = reqBody.Bio
	}
	if reqBody.Location != nil {
		req.Location = reqBody.Location
	}

	// 调用用户服务
	resp, err := handler.GetUserClient().UpdateUser(ctx, req)
	if err != nil {
		common.HandleServiceError(c, "UpdateUser", err, MsgUpdateUserInfoFailed)
		return
	}

	common.SuccessResponseFunc(c, resp)
}
