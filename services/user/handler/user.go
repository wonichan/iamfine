package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	service "hupu/kitex_gen/user"
	"hupu/services/user/repository"
	"hupu/shared/constants"
	"hupu/shared/log"
	"hupu/shared/middleware"
	"hupu/shared/models"
	"hupu/shared/utils"
)

type UserHandler struct {
	db *repository.UserRepository
}

func NewUserHandler() service.UserService {
	return &UserHandler{
		db: repository.NewUserRepository(),
	}
}

// 用户注册
func (h *UserHandler) Register(ctx context.Context, req *service.RegisterRequest) (*service.RegisterResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Info("Register start")

	// 参数验证
	validationErrors := middleware.ValidateUserRegistration(ctx, req.Username, req.Password, req.Email, req.Phone)
	if len(validationErrors) > 0 {
		errorMsg := "参数验证失败: "
		for _, ve := range validationErrors {
			errorMsg += fmt.Sprintf("%s: %s; ", ve.Field, ve.Message)
		}
		return &service.RegisterResponse{
			Code:    constants.ValidationErrorCode,
			Message: errorMsg,
		}, nil
	}

	userModel := &models.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Username,
		Phone:    req.Phone,
		Email:    req.Email,
		Status:   models.UserStatusActive,
	}

	savedUser, err := h.db.CreateUser(ctx, userModel)
	if err != nil {
		// 根据错误类型返回不同的错误码
		if strings.Contains(err.Error(), "already exists") {
			return &service.RegisterResponse{
				Code:    constants.UserAlreadyExistsCode,
				Message: constants.GetErrorMessage(constants.UserAlreadyExistsCode),
			}, nil
		}
		return &service.RegisterResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("failed to create user, err:%s", err.Error()),
		}, nil
	}
	logger.Infof("Register success, user:%+v", savedUser)
	return &service.RegisterResponse{
		Code: constants.SuccessCode,
		User: h.convertToUserResponse(savedUser),
	}, nil
}

// 用户登录
func (h *UserHandler) Login(ctx context.Context, req *service.LoginRequest) (*service.LoginResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Info("Login start")

	// 参数验证
	if req.Username == "" || req.Password == "" {
		return &service.LoginResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}

	getUser, err := h.db.GetUser(ctx, &models.User{
		Username: req.Username,
		Status:   models.UserStatusActive,
	})
	if err != nil {
		return &service.LoginResponse{
			Code:    constants.UserNotFoundCode,
			Message: constants.GetErrorMessage(constants.UserNotFoundCode),
		}, nil
	}
	if !utils.CheckPassword(req.Password, getUser.Password) {
		return &service.LoginResponse{
			Code:    constants.UserInvalidCredentialsCode,
			Message: constants.GetErrorMessage(constants.UserInvalidCredentialsCode),
		}, nil
	}

	return &service.LoginResponse{
		Code: constants.SuccessCode,
		User: h.convertToUserResponse(getUser),
	}, nil
}

// 获取用户信息
func (h *UserHandler) GetUser(ctx context.Context, req *service.GetUserRequest) (*service.GetUserResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Infof("GetUser start, req:%+v", req)

	// 参数验证
	if req.UserId == "" {
		return &service.GetUserResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}

	getUser, err := h.db.GetUser(ctx, &models.User{
		ID: req.UserId,
	})
	if err != nil {
		return &service.GetUserResponse{
			Code:    constants.UserNotFoundCode,
			Message: constants.GetErrorMessage(constants.UserNotFoundCode),
		}, nil
	}

	return &service.GetUserResponse{
		Code: constants.SuccessCode,
		User: h.convertToUserResponse(getUser),
	}, nil
}

// 更新用户信息
func (h *UserHandler) UpdateUser(ctx context.Context, req *service.UpdateUserRequest) (*service.UpdateUserResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Infof("UpdateUser start, req:%+v", req)

	// 参数验证
	if len(req.Id) <= 0 {
		return &service.UpdateUserResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}

	_, err := h.db.GetUser(ctx, &models.User{
		ID: req.Id,
	})
	if err != nil {
		return &service.UpdateUserResponse{
			Code:    constants.UserNotFoundCode,
			Message: constants.GetErrorMessage(constants.UserNotFoundCode),
		}, nil
	}

	// 处理可选更新字段
	newData := h.updateUserFields(req)
	newData.ID = req.Id

	updatedUser, err := h.db.UpdateUser(ctx, newData)
	if err != nil {
		return &service.UpdateUserResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("failed to update user, err:%s", err.Error()),
		}, nil
	}

	return &service.UpdateUserResponse{
		Code: constants.SuccessCode,
		User: h.convertToUserResponse(updatedUser),
	}, nil
}

// 关注功能相关方法
func (h *UserHandler) FollowUser(ctx context.Context, req *service.FollowUserRequest) (*service.FollowUserResponse, error) {
	// 参数验证
	validationErrors := middleware.ValidateFollowOperation(ctx, req.UserId, req.TargetUserId)
	if len(validationErrors) > 0 {
		errorMsg := "参数验证失败: "
		for _, ve := range validationErrors {
			errorMsg += fmt.Sprintf("%s: %s; ", ve.Field, ve.Message)
		}
		return &service.FollowUserResponse{
			Code:    constants.ValidationErrorCode,
			Message: errorMsg,
		}, nil
	}

	err := h.db.FollowUser(ctx, req.UserId, req.TargetUserId)
	if err != nil {
		// 根据错误类型返回不同的错误码
		if strings.Contains(err.Error(), "target user not found") {
			return &service.FollowUserResponse{
				Code:    constants.UserNotFoundCode,
				Message: constants.GetErrorMessage(constants.UserNotFoundCode),
			}, nil
		}
		if strings.Contains(err.Error(), "already following") {
			return &service.FollowUserResponse{
				Code:    constants.FollowAlreadyExistsCode,
				Message: constants.GetErrorMessage(constants.FollowAlreadyExistsCode),
			}, nil
		}
		return &service.FollowUserResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("failed to follow user, err:%s", err.Error()),
		}, nil
	}

	return &service.FollowUserResponse{
		Code: constants.SuccessCode,
	}, nil
}

func (h *UserHandler) UnfollowUser(ctx context.Context, req *service.UnfollowUserRequest) (*service.UnfollowUserResponse, error) {
	// 参数验证
	validationErrors := middleware.ValidateFollowOperation(ctx, req.UserId, req.TargetUserId)
	if len(validationErrors) > 0 {
		errorMsg := "参数验证失败: "
		for _, ve := range validationErrors {
			errorMsg += fmt.Sprintf("%s: %s; ", ve.Field, ve.Message)
		}
		return &service.UnfollowUserResponse{
			Code:    constants.ValidationErrorCode,
			Message: errorMsg,
		}, nil
	}

	err := h.db.UnfollowUser(ctx, req.UserId, req.TargetUserId)
	if err != nil {
		return &service.UnfollowUserResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("failed to unfollow user, err:%s", err.Error()),
		}, nil
	}

	return &service.UnfollowUserResponse{
		Code: constants.SuccessCode,
	}, nil
}

// 获取粉丝列表
func (h *UserHandler) GetFollowers(ctx context.Context, req *service.GetFollowersRequest) (*service.GetFollowersResponse, error) {
	// 参数验证
	if req.UserId == "" {
		return &service.GetFollowersResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}

	followers, err := h.db.GetFollowerList(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &service.GetFollowersResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("failed to get followers, err:%s", err.Error()),
		}, nil
	}

	// 转换为响应格式
	userResponses := make([]*service.User, len(followers))
	for i, follower := range followers {
		userResponses[i] = h.convertToUserResponse(follower)
	}

	return &service.GetFollowersResponse{
		Code:  constants.SuccessCode,
		Users: userResponses,
		Total: int32(len(followers)),
	}, nil
}

// 获取关注列表
func (h *UserHandler) GetFollowing(ctx context.Context, req *service.GetFollowingRequest) (*service.GetFollowingResponse, error) {
	// 参数验证
	if req.UserId == "" {
		return &service.GetFollowingResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}

	following, err := h.db.GetFollowingList(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &service.GetFollowingResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("failed to get following list, err:%s", err.Error()),
		}, nil
	}

	// 转换为响应格式
	userResponses := make([]*service.User, len(following))
	for i, follow := range following {
		userResponses[i] = h.convertToUserResponse(follow)
	}

	return &service.GetFollowingResponse{
		Code:  constants.SuccessCode,
		Users: userResponses,
		Total: int32(len(following)),
	}, nil
}

// 匿名马甲管理相关方法
func (h *UserHandler) CreateAnonymousProfile(ctx context.Context, req *service.CreateAnonymousProfileRequest) (*service.CreateAnonymousProfileResponse, error) {
	// 参数验证
	if req.UserId == "" || req.AnonymousName == "" {
		return &service.CreateAnonymousProfileResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}

	avatar := &models.AnonymousAvatar{
		UserID:   req.UserId,
		Name:     req.AnonymousName,
		Color:    req.AvatarColor,
		IsActive: true,
	}

	err := h.db.CreateAnonymousAvatar(ctx, avatar)
	if err != nil {
		return &service.CreateAnonymousProfileResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("failed to create anonymous avatar, err:%s", err.Error()),
		}, nil
	}

	return &service.CreateAnonymousProfileResponse{
		Code: constants.SuccessCode,
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
	// 参数验证
	if req.UserId == "" {
		return &service.GetAnonymousProfilesResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}

	avatars, err := h.db.GetAnonymousAvatarList(ctx, req.UserId)
	if err != nil {
		return &service.GetAnonymousProfilesResponse{
			Code:    constants.DatabaseErrorCode,
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
		Code:     constants.SuccessCode,
		Profiles: avatarResponses,
	}, nil
}

func (h *UserHandler) UpdateAnonymousProfile(ctx context.Context, req *service.UpdateAnonymousProfileRequest) (*service.UpdateAnonymousProfileResponse, error) {
	// 参数验证
	if req.ProfileId == "" {
		return &service.UpdateAnonymousProfileResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}

	// 先获取现有的匿名头像
	avatar, err := h.db.GetAnonymousAvatar(ctx, req.ProfileId)
	if err != nil {
		return &service.UpdateAnonymousProfileResponse{
			Code:    constants.DatabaseErrorCode,
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
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("failed to update anonymous avatar, err:%s", err.Error()),
		}, nil
	}

	return &service.UpdateAnonymousProfileResponse{
		Code: constants.SuccessCode,
	}, nil
}

// 用户统计相关方法
func (h *UserHandler) GetUserStats(ctx context.Context, req *service.GetUserStatsRequest) (*service.GetUserStatsResponse, error) {
	// 参数验证
	if req.UserId == "" {
		return &service.GetUserStatsResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}

	stats, err := h.db.GetUserStats(ctx, req.UserId)
	if err != nil {
		return &service.GetUserStatsResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("failed to get user stats, err:%s", err.Error()),
		}, nil
	}

	return &service.GetUserStatsResponse{
		Code:           constants.SuccessCode,
		PostCount:      stats.PostCount,
		CommentCount:   stats.CommentCount,
		LikeCount:      stats.LikeCount,
		CollectCount:   stats.FavoriteCount,
		AverageScore:   stats.AverageScore,
		FollowerCount:  stats.FollowerCount,
		FollowingCount: stats.FollowingCount,
	}, nil
}

// 辅助方法：更新用户字段
func (h *UserHandler) updateUserFields(req *service.UpdateUserRequest) *models.User {
	updateData := &models.User{}
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
	if req.Email != nil {
		updateData.Email = *req.Email
	}
	if req.Phone != nil {
		updateData.Phone = *req.Phone
	}
	return updateData
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
