package handler

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	service "hupu/kitex_gen/user"
	"hupu/services/user/repository"
	"hupu/shared/constants"
	"hupu/shared/log"
	"hupu/shared/models"
	"hupu/shared/utils"
)

type UserHandler struct {
	db repository.UserRepository
}

func NewUserHandler(db *gorm.DB) service.UserService {
	return &UserHandler{
		db: repository.NewUserRepository(db),
	}
}

// 用户注册
func (h *UserHandler) Register(ctx context.Context, req *service.RegisterRequest) (*service.RegisterResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Info("Register start")
	userModel := &models.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Username,
		Phone:    req.Phone,
		Status:   models.UserStatusActive,
	}

	savedUser, err := h.db.CreateUser(ctx, userModel)
	if err != nil {
		return &service.RegisterResponse{
			Code:    constants.UserExistsErrCode,
			Message: fmt.Sprintf("failed to create user, err:%s", err.Error()),
		}, nil
	}
	logger.Infof("Register success, user:%+v", savedUser)
	return &service.RegisterResponse{
		Code: 0,
		User: h.convertToUserResponse(savedUser),
	}, nil
}

// 用户登录
func (h *UserHandler) Login(ctx context.Context, req *service.LoginRequest) (*service.LoginResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Info("Login start")
	getUser, err := h.db.GetUser(ctx, &models.User{
		Username: req.Username,
		Password: req.Password,
		Status:   models.UserStatusActive,
	})
	if err != nil {
		return &service.LoginResponse{
			Code:    constants.UserLoginErrCode,
			Message: fmt.Sprintf("failed to login, err:%s", err.Error()),
		}, nil
	}
	if !utils.CheckPassword(req.Password, getUser.Password) {
		return &service.LoginResponse{
			Code:    constants.UserPasswdErrCode,
			Message: fmt.Sprintf("failed to login, err:%s", constants.MsgPasswordError),
		}, nil
	}

	return &service.LoginResponse{
		Code: 0,
		User: h.convertToUserResponse(getUser),
	}, nil
}

// 获取用户信息
func (h *UserHandler) GetUser(ctx context.Context, req *service.GetUserRequest) (*service.GetUserResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Infof("GetUser start, req:%+v", req)
	getUser, err := h.db.GetUser(ctx, &models.User{
		ID: req.UserId,
	})
	if err != nil {
		return &service.GetUserResponse{
			Code:    constants.UserGerUserErrCode,
			Message: fmt.Sprintf("failed to get user, err:%s", err.Error()),
		}, nil
	}

	return &service.GetUserResponse{
		Code: 0,
		User: h.convertToUserResponse(getUser),
	}, nil
}

// 更新用户信息
func (h *UserHandler) UpdateUser(ctx context.Context, req *service.UpdateUserRequest) (*service.UpdateUserResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Infof("UpdateUser start, req:%+v", req)
	updateData, err := h.db.GetUser(ctx, &models.User{
		ID: req.Id,
	})
	if err != nil {
		return &service.UpdateUserResponse{
			Code:    constants.UserUpdateUserErrCode,
			Message: fmt.Sprintf("failed to get user, err:%s", err.Error()),
		}, nil
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
		status := service.RelationshipStatus(*req.RelationshipStatus)
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
	updatedUser, err := h.db.UpdateUser(ctx, updateData)
	if err != nil {
		return &service.UpdateUserResponse{
			Code:    constants.UserUpdateUserErrCode,
			Message: fmt.Sprintf("failed to update user, err:%s", err.Error()),
		}, nil
	}

	return &service.UpdateUserResponse{
		Code: 0,
		User: h.convertToUserResponse(updatedUser),
	}, nil
}

// 关注功能相关方法
func (h *UserHandler) FollowUser(ctx context.Context, req *service.FollowUserRequest) (*service.FollowUserResponse, error) {
	// 不能关注自己
	if req.UserId == req.TargetUserId {
		return &service.FollowUserResponse{
			Code:    constants.UserFollowUserErrCode,
			Message: constants.MsgCannotFollowSelf,
		}, nil
	}
	err := h.db.FollowUser(ctx, req.UserId, req.TargetUserId)
	if err != nil {
		return &service.FollowUserResponse{
			Code:    constants.UserFollowUserErrCode,
			Message: fmt.Sprintf("failed to follow user, err:%s", err.Error()),
		}, nil
	}

	return &service.FollowUserResponse{
		Code: 0,
	}, nil
}

func (h *UserHandler) UnfollowUser(ctx context.Context, req *service.UnfollowUserRequest) (*service.UnfollowUserResponse, error) {
	err := h.db.UnfollowUser(ctx, req.UserId, req.TargetUserId)
	if err != nil {
		return &service.UnfollowUserResponse{
			Code:    constants.UserUnfollowUserErrCode,
			Message: fmt.Sprintf("failed to unfollow user, err:%s", err.Error()),
		}, nil
	}

	return &service.UnfollowUserResponse{
		Code: 0,
	}, nil
}

// 获取粉丝列表
func (h *UserHandler) GetFollowers(ctx context.Context, req *service.GetFollowersRequest) (*service.GetFollowersResponse, error) {
	followers, err := h.db.GetFollowerList(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &service.GetFollowersResponse{
			Code:    constants.UserGetFollowersErrCode,
			Message: fmt.Sprintf("failed to get followers, err:%s", err.Error()),
		}, nil
	}

	// 转换为响应格式
	userResponses := make([]*service.User, len(followers))
	for i, follower := range followers {
		userResponses[i] = h.convertToUserResponse(follower)
	}

	return &service.GetFollowersResponse{
		Code:  0,
		Users: userResponses,
		Total: int32(len(followers)),
	}, nil
}

// 获取关注列表
func (h *UserHandler) GetFollowing(ctx context.Context, req *service.GetFollowingRequest) (*service.GetFollowingResponse, error) {
	following, err := h.db.GetFollowingList(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &service.GetFollowingResponse{
			Code:    constants.UserGetFollowingErrCode,
			Message: fmt.Sprintf("failed to get following list, err:%s", err.Error()),
		}, nil
	}

	// 转换为响应格式
	userResponses := make([]*service.User, len(following))
	for i, follow := range following {
		userResponses[i] = h.convertToUserResponse(follow)
	}

	return &service.GetFollowingResponse{
		Code:  0,
		Users: userResponses,
		Total: int32(len(following)),
	}, nil
}

// 匿名马甲管理相关方法
func (h *UserHandler) CreateAnonymousProfile(ctx context.Context, req *service.CreateAnonymousProfileRequest) (*service.CreateAnonymousProfileResponse, error) {
	avatar := &models.AnonymousAvatar{
		UserID:   req.UserId,
		Name:     req.AnonymousName,
		Color:    req.AvatarColor,
		IsActive: true,
	}

	err := h.db.CreateAnonymousAvatar(ctx, avatar)
	if err != nil {
		return &service.CreateAnonymousProfileResponse{
			Code:    constants.UserCreateAnonymousErrCode,
			Message: fmt.Sprintf("failed to create anonymous avatar, err:%s", err.Error()),
		}, nil
	}

	return &service.CreateAnonymousProfileResponse{
		Code: 0,
		Profile: &service.AnonymousProfile{
			Id:            avatar.ID,
			UserId:        avatar.UserID,
			AnonymousName: avatar.Name,
			AvatarColor:   avatar.Color,
			IsActive:      avatar.IsActive,
			CreatedAt:     time.Now().Unix(),
		},
	}, nil
}

func (h *UserHandler) GetAnonymousProfiles(ctx context.Context, req *service.GetAnonymousProfilesRequest) (*service.GetAnonymousProfilesResponse, error) {
	avatars, err := h.db.GetAnonymousAvatarList(ctx, req.UserId)
	if err != nil {
		return &service.GetAnonymousProfilesResponse{
			Code:    constants.UserGetAnonymousErrCode,
			Message: fmt.Sprintf("failed to get anonymous avatar list, err:%s", err.Error()),
		}, nil
	}

	// 转换为响应格式
	avatarResponses := make([]*service.AnonymousProfile, len(avatars))
	for i, avatar := range avatars {
		avatarResponses[i] = &service.AnonymousProfile{
			Id:            avatar.ID,
			UserId:        avatar.UserID,
			AnonymousName: avatar.Name,
			AvatarColor:   avatar.Color,
			IsActive:      avatar.IsActive,
			CreatedAt:     time.Now().Unix(),
		}
	}

	return &service.GetAnonymousProfilesResponse{
		Code:     0,
		Profiles: avatarResponses,
	}, nil
}

func (h *UserHandler) UpdateAnonymousProfile(ctx context.Context, req *service.UpdateAnonymousProfileRequest) (*service.UpdateAnonymousProfileResponse, error) {
	// 先获取现有的匿名头像
	avatar, err := h.db.GetAnonymousAvatar(ctx, req.ProfileId)
	if err != nil {
		return &service.UpdateAnonymousProfileResponse{
			Code:    constants.UserGetAnonymousErrCode,
			Message: fmt.Sprintf("failed to get anonymous avatar, err:%s", err.Error()),
		}, nil
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
		return &service.UpdateAnonymousProfileResponse{
			Code:    constants.UserUpdateAnonymousErrCode,
			Message: fmt.Sprintf("failed to update anonymous avatar, err:%s", err.Error()),
		}, nil
	}

	return &service.UpdateAnonymousProfileResponse{
		Code: 0,
	}, nil
}

// 用户统计相关方法
func (h *UserHandler) GetUserStats(ctx context.Context, req *service.GetUserStatsRequest) (*service.GetUserStatsResponse, error) {
	stats, err := h.db.GetUserStats(ctx, req.UserId)
	if err != nil {
		return &service.GetUserStatsResponse{
			Code:    constants.UserGetUserStatusErrCode,
			Message: fmt.Sprintf("failed to get user stats, err:%s", err.Error()),
		}, nil
	}

	return &service.GetUserStatsResponse{
		Code:           0,
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
func (h *UserHandler) convertToUserResponse(u *models.User) *service.User {
	response := &service.User{
		Id:             u.ID,
		Username:       u.Username,
		Nickname:       u.Nickname,
		Avatar:         u.Avatar,
		Phone:          u.Phone,
		Email:          u.Email,
		Status:         service.UserStatus(u.Status),
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
		status := service.RelationshipStatus(*u.RelationshipStatus)
		response.RelationshipStatus = status
	}
	if u.AgeGroup != nil {
		age := service.AgeGroup(*u.AgeGroup)
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
