package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"

	"hupu/kitex_gen/user"
	"hupu/shared/log"
)

// 用户注册
func Register(ctx context.Context, c *app.RequestContext) {
	logger.Info("Register start")
	// 解析请求参数
	var req user.RegisterRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 调用用户服务
	resp, err := userClient.Register(ctx, &req)
	if err != nil {
		log.GetLogger().Errorf("Register error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "注册失败",
		})
		return
	}
	logger.Infof("Register success, resp:%+v", resp)
	c.JSON(http.StatusOK, resp)
}

// 用户登录
func Login(ctx context.Context, c *app.RequestContext) {
	// 解析请求参数
	var req user.LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	// 调用用户服务
	resp, err := userClient.Login(ctx, &req)
	if err != nil {
		log.GetLogger().Errorf("Login error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "登录失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取用户信息
func GetUser(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 调用用户服务
	req := &user.GetUserRequest{
		UserId: userID,
	}
	resp, err := userClient.GetUser(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetUser error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 更新用户信息
func UpdateUser(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取要更新的用户ID
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 检查权限：只能更新自己的信息
	if currentUserID.(string) != userID {
		c.JSON(http.StatusForbidden, map[string]interface{}{
			"code":    403,
			"message": "无权限更新其他用户信息",
		})
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
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数格式错误",
		})
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
	resp, err := userClient.UpdateUser(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("UpdateUser error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新用户信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 关注用户
func FollowUser(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取要关注的用户ID
	targetUserID := c.Param("id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "目标用户ID不能为空",
		})
		return
	}

	// 不能关注自己
	if userID.(string) == targetUserID {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "不能关注自己",
		})
		return
	}

	// 构建请求
	req := &user.FollowUserRequest{
		UserId:       userID.(string),
		TargetUserId: targetUserID,
	}

	// 调用用户服务
	resp, err := userClient.FollowUser(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("FollowUser error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "关注用户失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 取消关注用户
func UnfollowUser(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取要取消关注的用户ID
	targetUserID := c.Param("id")
	if targetUserID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "目标用户ID不能为空",
		})
		return
	}

	// 构建请求
	req := &user.UnfollowUserRequest{
		UserId:       userID.(string),
		TargetUserId: targetUserID,
	}

	// 调用用户服务
	resp, err := userClient.UnfollowUser(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("UnfollowUser error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "取消关注用户失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取粉丝列表
func GetFollowers(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 解析查询参数
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 构建请求
	req := &user.GetFollowersRequest{
		UserId:   userID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	// 调用用户服务
	resp, err := userClient.GetFollowers(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetFollowers error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取粉丝列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取关注列表
func GetFollowing(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 解析查询参数
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 构建请求
	req := &user.GetFollowingRequest{
		UserId:   userID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	// 调用用户服务
	resp, err := userClient.GetFollowing(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetFollowing error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取关注列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 创建匿名马甲
func CreateAnonymousProfile(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 解析请求体
	var reqBody struct {
		AnonymousName string `json:"anonymous_name" binding:"required"`
		AvatarColor   string `json:"avatar_color"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数格式错误",
		})
		return
	}

	// 构建请求
	req := &user.CreateAnonymousProfileRequest{
		UserId:        userID.(string),
		AnonymousName: reqBody.AnonymousName,
		AvatarColor:   reqBody.AvatarColor,
	}

	// 调用用户服务
	resp, err := userClient.CreateAnonymousProfile(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("CreateAnonymousProfile error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建匿名马甲失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取匿名马甲列表
func GetAnonymousProfiles(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 构建请求
	req := &user.GetAnonymousProfilesRequest{
		UserId: userID.(string),
	}

	// 调用用户服务
	resp, err := userClient.GetAnonymousProfiles(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetAnonymousProfiles error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取匿名马甲列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 更新匿名马甲
func UpdateAnonymousProfile(ctx context.Context, c *app.RequestContext) {

	// 获取profile_id
	profileID := c.Param("profile_id")
	if profileID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "profile_id不能为空",
		})
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
	resp, err := userClient.UpdateAnonymousProfile(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("UpdateAnonymousProfile error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新匿名马甲失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取用户统计信息
func GetUserStats(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID参数
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 构建请求
	req := &user.GetUserStatsRequest{
		UserId: userID,
	}

	// 调用用户服务
	resp, err := userClient.GetUserStats(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetUserStats error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取用户统计信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
