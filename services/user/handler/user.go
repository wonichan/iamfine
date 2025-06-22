package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"hupu/kitex_gen/user"
	"hupu/services/user/repository"
	"hupu/shared/constants"
	"hupu/shared/log"
	"hupu/shared/models"
	"hupu/shared/utils"
)

type UserHandler struct {
	db  repository.UserRepository
	rdb repository.UserRepository
}

func NewUserHandler(db *gorm.DB, rdb *redis.Client) *UserHandler {
	return &UserHandler{
		db:  repository.NewUserRepository(db),
		rdb: repository.NewUserRedisRepo(rdb),
	}
}

// 用户注册
func (h *UserHandler) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Info("Register start")
	userModel := &models.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Username,
		Phone:    req.Phone,
		Status:   models.UserStatusActive,
	}

	savedUser, err := h.rdb.CreateUser(ctx, userModel)
	if err != nil {
		return &user.RegisterResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to create user, err:%s", err.Error()),
		}, nil
	}
	logger.Infof("Register success, user:%+v", savedUser)
	return &user.RegisterResponse{
		Code:    200,
		Message: "注册成功",
		User:    h.convertToUserResponse(savedUser),
	}, nil
}

// 用户登录
func (h *UserHandler) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Info("Login start")
	getUser, err := h.rdb.GetUser(ctx, &models.User{
		Username: req.Username,
		Password: req.Password,
		Status:   models.UserStatusActive,
	})
	if err != nil {
		return &user.LoginResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to login, err:%s", err.Error()),
		}, nil
	}

	// 生成JWT token
	token, err := utils.GenerateToken(getUser.ID, getUser.Username)
	if err != nil {
		return &user.LoginResponse{
			Code:    500,
			Message: "生成token失败",
		}, err
	}

	return &user.LoginResponse{
		Code:    200,
		Message: "登录成功",
		Token:   token,
		User:    h.convertToUserResponse(getUser),
	}, nil
}

// 获取用户信息
func (h *UserHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Infof("GetUser start, req:%+v", req)
	getUser, err := h.rdb.GetUser(ctx, &models.User{
		ID: req.UserId,
	})
	if err != nil {
		return &user.GetUserResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to get user, err:%s", err.Error()),
		}, err
	}

	return &user.GetUserResponse{
		Code:    0,
		Message: "查询成功",
		User:    h.convertToUserResponse(getUser),
	}, nil
}

// 更新用户信息
func (h *UserHandler) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Infof("UpdateUser start, req:%+v", req)

	updateData := &models.User{
		ID: req.Id,
	}

	// 处理可选更新字段
	if req.Nickname != nil {
		updateData.Nickname = *req.Nickname
	}
	if req.Avatar != nil {
		updateData.Avatar = *req.Avatar
	}
	if req.Bio != nil {
		updateData.Bio = req.Bio
	}
	if req.RelationshipStatus != nil {
		status := models.RelationshipStatus(*req.RelationshipStatus)
		updateData.RelationshipStatus = &status
	}
	if req.AgeGroup != nil {
		age := models.AgeGroup(*req.AgeGroup)
		updateData.AgeGroup = &age
	}
	if req.Location != nil {
		updateData.Location = req.Location
	}
	if req.Tags != nil {
		updateData.Tags = models.StringArray(req.Tags)
	}

	updatedUser, err := h.rdb.UpdateUser(ctx, updateData)
	if err != nil {
		return &user.UpdateUserResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to update user, err:%s", err.Error()),
		}, err
	}

	return &user.UpdateUserResponse{
		Code:    0,
		Message: "更新成功",
		User:    h.convertToUserResponse(updatedUser),
	}, nil
}

// 关注功能相关方法
func (h *UserHandler) FollowUser(ctx context.Context, req *user.FollowUserRequest) (*user.FollowUserResponse, error) {
	err := h.rdb.FollowUser(ctx, req.UserId, req.TargetUserId)
	if err != nil {
		return &user.FollowUserResponse{
			Code:    500,
			Message: "关注失败",
		}, err
	}

	return &user.FollowUserResponse{
		Code:    0,
		Message: "关注成功",
	}, nil
}

func (h *UserHandler) UnfollowUser(ctx context.Context, req *user.UnfollowUserRequest) (*user.UnfollowUserResponse, error) {
	err := h.rdb.UnfollowUser(ctx, req.UserId, req.TargetUserId)
	if err != nil {
		return &user.UnfollowUserResponse{
			Code:    500,
			Message: "取消关注失败",
		}, err
	}

	return &user.UnfollowUserResponse{
		Code:    0,
		Message: "取消关注成功",
	}, nil
}

// 获取粉丝列表
func (h *UserHandler) GetFollowers(ctx context.Context, req *user.GetFollowersRequest) (*user.GetFollowersResponse, error) {
	followers, err := h.db.GetFollowerList(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &user.GetFollowersResponse{
			Code:    500,
			Message: "获取粉丝列表失败",
		}, err
	}

	// 转换为响应格式
	userResponses := make([]*user.User, len(followers))
	for i, follower := range followers {
		userResponses[i] = h.convertToUserResponse(follower)
	}

	return &user.GetFollowersResponse{
		Code:    200,
		Message: "获取成功",
		Users:   userResponses,
		Total:   int32(len(followers)),
	}, nil
}

// 获取关注列表
func (h *UserHandler) GetFollowing(ctx context.Context, req *user.GetFollowingRequest) (*user.GetFollowingResponse, error) {
	following, err := h.db.GetFollowingList(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &user.GetFollowingResponse{
			Code:    500,
			Message: "获取关注列表失败",
		}, err
	}

	// 转换为响应格式
	userResponses := make([]*user.User, len(following))
	for i, follow := range following {
		userResponses[i] = h.convertToUserResponse(follow)
	}

	return &user.GetFollowingResponse{
		Code:    200,
		Message: "获取成功",
		Users:   userResponses,
		Total:   int32(len(following)),
	}, nil
}

// 匿名马甲管理相关方法
func (h *UserHandler) CreateAnonymousProfile(ctx context.Context, req *user.CreateAnonymousProfileRequest) (*user.CreateAnonymousProfileResponse, error) {
	avatar := &models.AnonymousAvatar{
		UserID:   req.UserId,
		Name:     req.AnonymousName,
		Color:    req.AvatarColor,
		IsActive: true,
	}

	err := h.db.CreateAnonymousAvatar(ctx, avatar)
	if err != nil {
		return &user.CreateAnonymousProfileResponse{
			Code:    500,
			Message: "创建匿名马甲失败",
		}, err
	}

	return &user.CreateAnonymousProfileResponse{
		Code:    200,
		Message: "创建成功",
		Profile: &user.AnonymousProfile{
			Id:            avatar.ID,
			UserId:        avatar.UserID,
			AnonymousName: avatar.Name,
			AvatarColor:   avatar.Color,
			IsActive:      avatar.IsActive,
			CreatedAt:     time.Now().Unix(),
		},
	}, nil
}

func (h *UserHandler) GetAnonymousProfiles(ctx context.Context, req *user.GetAnonymousProfilesRequest) (*user.GetAnonymousProfilesResponse, error) {
	avatars, err := h.db.GetAnonymousAvatarList(ctx, req.UserId)
	if err != nil {
		return &user.GetAnonymousProfilesResponse{
			Code:    500,
			Message: "获取匿名马甲列表失败",
		}, err
	}

	// 转换为响应格式
	avatarResponses := make([]*user.AnonymousProfile, len(avatars))
	for i, avatar := range avatars {
		avatarResponses[i] = &user.AnonymousProfile{
			Id:            avatar.ID,
			UserId:        avatar.UserID,
			AnonymousName: avatar.Name,
			AvatarColor:   avatar.Color,
			IsActive:      avatar.IsActive,
			CreatedAt:     time.Now().Unix(),
		}
	}

	return &user.GetAnonymousProfilesResponse{
		Code:     200,
		Message:  "获取成功",
		Profiles: avatarResponses,
	}, nil
}

func (h *UserHandler) UpdateAnonymousProfile(ctx context.Context, req *user.UpdateAnonymousProfileRequest) (*user.UpdateAnonymousProfileResponse, error) {
	// 先获取现有的匿名头像
	avatar, err := h.db.GetAnonymousAvatar(ctx, req.ProfileId)
	if err != nil {
		return &user.UpdateAnonymousProfileResponse{
			Code:    500,
			Message: "获取匿名马甲失败",
		}, err
	}

	// 更新字段
	if req.AnonymousName != nil {
		avatar.Name = *req.AnonymousName
	}
	if req.AvatarColor != nil {
		avatar.Color = *req.AvatarColor
	}
	if req.IsActive != nil {
		avatar.IsActive = *req.IsActive
	}

	err = h.db.UpdateAnonymousAvatar(ctx, avatar)
	if err != nil {
		return &user.UpdateAnonymousProfileResponse{
			Code:    500,
			Message: "更新匿名马甲失败",
		}, err
	}

	return &user.UpdateAnonymousProfileResponse{
		Code:    200,
		Message: "更新成功",
	}, nil
}

// 用户统计相关方法
func (h *UserHandler) GetUserStats(ctx context.Context, req *user.GetUserStatsRequest) (*user.GetUserStatsResponse, error) {
	stats, err := h.db.GetUserStats(ctx, req.UserId)
	if err != nil {
		return &user.GetUserStatsResponse{
			Code:    500,
			Message: "获取用户统计失败",
		}, err
	}

	return &user.GetUserStatsResponse{
		Code:           200,
		Message:        "获取成功",
		PostCount:      stats.PostCount,
		CommentCount:   stats.CommentCount,
		LikeCount:      stats.LikeCount,
		CollectCount:   stats.FavoriteCount,
		AverageScore:   stats.AverageScore,
		FollowerCount:  stats.FollowerCount,
		FollowingCount: stats.FollowingCount,
	}, nil
}

// 辅助方法：转换模型为响应格式
func (h *UserHandler) convertToUserResponse(u *models.User) *user.User {
	response := &user.User{
		Id:             u.ID,
		Username:       u.Username,
		Nickname:       u.Nickname,
		Avatar:         u.Avatar,
		Phone:          u.Phone,
		Email:          u.Email,
		Status:         user.UserStatus(u.Status),
		PostCount:      u.PostCount,
		CommentCount:   u.CommentCount,
		LikeCount:      u.LikeCount,
		CollectCount:   u.FavoriteCount,
		AverageScore:   u.AverageScore,
		FollowerCount:  u.FollowerCount,
		FollowingCount: u.FollowingCount,
		IsVerified:     u.IsVerified,
		CreatedAt:      u.CreatedAt.Unix(),
		UpdatedAt:      u.UpdatedAt.Unix(),
	}

	if u.Bio != nil {
		response.Bio = u.Bio
	}
	if u.RelationshipStatus != nil {
		status := user.RelationshipStatus(*u.RelationshipStatus)
		response.RelationshipStatus = status
	}
	if u.AgeGroup != nil {
		age := user.AgeGroup(*u.AgeGroup)
		response.AgeGroup = age
	}
	if u.Location != nil {
		response.Location = u.Location
	}
	if len(u.Tags) > 0 {
		tags := make([]string, len(u.Tags))
		copy(tags, u.Tags)
		response.Tags = tags
	}

	return response
}
