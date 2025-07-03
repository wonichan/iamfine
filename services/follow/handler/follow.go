package handler

import (
	"context"
	"fmt"
	"hupu/kitex_gen/follow"
	"hupu/services/follow/repository"
	"hupu/shared/constants"
	"hupu/shared/middleware"

	"gorm.io/gorm"
)

type FollowHandler struct {
	db repository.FollowRepository
}

func NewFollowHandler(db *gorm.DB) *FollowHandler {
	return &FollowHandler{
		db: repository.NewFollowRepository(db),
	}
}

func (h *FollowHandler) Follow(ctx context.Context, req *follow.FollowRequest) (*follow.FollowResponse, error) {
	// 参数验证
	if err := middleware.ValidateFollowOperation(req.FollowerId, req.FollowingId); err != nil {
		return &follow.FollowResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}
	
	err := h.db.Follow(ctx, req.FollowerId, req.FollowingId)
	if err != nil {
		if fmt.Sprintf("%s", err).Contains("already following") {
			return &follow.FollowResponse{
				Code:    constants.UserAlreadyFollowedCode,
				Message: constants.GetErrorMessage(constants.UserAlreadyFollowedCode),
			}, nil
		}
		return &follow.FollowResponse{
			Code:    constants.FollowOperationFailedCode,
			Message: fmt.Sprintf("关注失败: %s", err),
		}, nil
	}

	return &follow.FollowResponse{
		Code:    constants.SuccessCode,
		Message: "关注成功",
	}, nil
}

func (h *FollowHandler) Unfollow(ctx context.Context, req *follow.UnfollowRequest) (*follow.UnfollowResponse, error) {
	// 参数验证
	if err := middleware.ValidateFollowOperation(req.FollowerId, req.FollowingId); err != nil {
		return &follow.UnfollowResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}
	
	err := h.db.Unfollow(ctx, req.FollowerId, req.FollowingId)
	if err != nil {
		if fmt.Sprintf("%s", err).Contains("not following") {
			return &follow.UnfollowResponse{
				Code:    constants.UserNotFollowedCode,
				Message: constants.GetErrorMessage(constants.UserNotFollowedCode),
			}, nil
		}
		return &follow.UnfollowResponse{
			Code:    constants.UnfollowOperationFailedCode,
			Message: fmt.Sprintf("取消关注失败: %s", err),
		}, nil
	}

	return &follow.UnfollowResponse{
		Code:    constants.SuccessCode,
		Message: "取消关注成功",
	}, nil
}

func (h *FollowHandler) IsFollowing(ctx context.Context, req *follow.FollowRequest) (*follow.FollowResponse, error) {
	// 参数验证
	if err := middleware.ValidateFollowOperation(req.FollowerId, req.FollowingId); err != nil {
		return &follow.FollowResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}
	
	isFollowing, err := h.db.IsFollowing(ctx, req.FollowerId, req.FollowingId)
	if err != nil {
		return &follow.FollowResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("查询失败: %s", err),
		}, nil
	}

	if !isFollowing {
		return &follow.FollowResponse{
			Code:    constants.UserNotFollowedCode,
			Message: constants.GetErrorMessage(constants.UserNotFollowedCode),
		}, nil
	}

	return &follow.FollowResponse{
		Code:        constants.SuccessCode,
		Message:     "已关注",
		IsFollowing: true,
	}, nil
}

func (h *FollowHandler) GetFollowList(ctx context.Context, req *follow.GetFollowListRequest) (*follow.GetFollowListResponse, error) {
	// 参数验证
	if req.UserId <= 0 {
		return &follow.GetFollowListResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}
	
	followList, err := h.db.GetFollowList(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &follow.GetFollowListResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("查询关注列表失败: %s", err),
		}, nil
	}

	// 转换为Follow结构体数组
	follows := make([]*follow.Follow, len(followList))
	for i, userID := range followList {
		follows[i] = &follow.Follow{
			FollowerId:  req.UserId,
			FollowingId: userID,
		}
	}

	return &follow.GetFollowListResponse{
		Code:    constants.SuccessCode,
		Message: "查询成功",
		Follows: follows,
		Total:   int32(len(follows)),
		HasMore: len(follows) == int(req.PageSize),
	}, nil
}

func (h *FollowHandler) GetFollowerList(ctx context.Context, req *follow.GetFollowerListRequest) (*follow.GetFollowerListResponse, error) {
	// 参数验证
	if req.UserId <= 0 {
		return &follow.GetFollowerListResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}
	
	follows, err := h.db.GetFollowerList(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &follow.GetFollowerListResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("查询粉丝列表失败: %s", err),
		}, nil
	}

	// 转换为响应格式
	var followerList []*follow.Follow
	for _, f := range follows {
		followerList = append(followerList, &follow.Follow{
			FollowerId:  f.FollowerID,
			FollowingId: f.FollowingID,
			CreatedAt:   f.CreatedAt.Unix(),
		})
	}

	return &follow.GetFollowerListResponse{
		Code:      constants.SuccessCode,
		Message:   "查询成功",
		Followers: followerList,
	}, nil
}

// GetFollowCount 获取关注数量
func (h *FollowHandler) GetFollowCount(ctx context.Context, req *follow.GetFollowCountRequest) (*follow.GetFollowCountResponse, error) {
	// 参数验证
	if req.UserId <= 0 {
		return &follow.GetFollowCountResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}
	
	count, err := h.db.GetFollowCount(ctx, req.UserId)
	if err != nil {
		return &follow.GetFollowCountResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("查询关注数量失败: %s", err),
		}, nil
	}

	return &follow.GetFollowCountResponse{
		Code:    constants.SuccessCode,
		Message: "查询成功",
		Count:   int32(count),
	}, nil
}

// GetFollowerCount 获取粉丝数量
func (h *FollowHandler) GetFollowerCount(ctx context.Context, req *follow.GetFollowerCountRequest) (*follow.GetFollowerCountResponse, error) {
	// 参数验证
	if req.UserId <= 0 {
		return &follow.GetFollowerCountResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}
	
	count, err := h.db.GetFollowerCount(ctx, req.UserId)
	if err != nil {
		return &follow.GetFollowerCountResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("查询粉丝数量失败: %s", err),
		}, nil
	}

	return &follow.GetFollowerCountResponse{
		Code:    constants.SuccessCode,
		Message: "查询成功",
		Count:   int32(count),
	}, nil
}

// GetMutualFollows 获取共同关注
func (h *FollowHandler) GetMutualFollows(ctx context.Context, req *follow.GetMutualFollowsRequest) (*follow.GetMutualFollowsResponse, error) {
	// 参数验证
	if req.UserId <= 0 || req.TargetUserId <= 0 {
		return &follow.GetMutualFollowsResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}
	
	userIds, err := h.db.GetMutualFollows(ctx, req.UserId, req.TargetUserId, req.Page, req.PageSize)
	if err != nil {
		return &follow.GetMutualFollowsResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("查询共同关注失败: %s", err),
		}, nil
	}

	// 转换为Follow结构体数组
	mutualFollows := make([]*follow.Follow, len(userIds))
	for i, userID := range userIds {
		mutualFollows[i] = &follow.Follow{
			FollowerId:  req.UserId,
			FollowingId: userID,
			IsMutual:    true,
		}
	}

	return &follow.GetMutualFollowsResponse{
		Code:          constants.SuccessCode,
		Message:       "查询成功",
		MutualFollows: mutualFollows,
		Total:         int32(len(mutualFollows)),
		HasMore:       len(mutualFollows) == int(req.PageSize),
	}, nil
}

// CheckFollowStatus 检查关注状态
func (h *FollowHandler) CheckFollowStatus(ctx context.Context, req *follow.CheckFollowStatusRequest) (*follow.CheckFollowStatusResponse, error) {
	// 参数验证
	if err := middleware.ValidateFollowOperation(req.FollowerId, req.FollowingId); err != nil {
		return &follow.CheckFollowStatusResponse{
			Code:    constants.ValidationErrorCode,
			Message: constants.GetErrorMessage(constants.ValidationErrorCode),
		}, nil
	}
	
	// 检查是否关注
	isFollowing, err := h.db.IsFollowing(ctx, req.FollowerId, req.FollowingId)
	if err != nil {
		return &follow.CheckFollowStatusResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("查询关注状态失败: %s", err),
		}, nil
	}

	// 如果已关注，检查是否互相关注
	var isMutual bool
	if isFollowing {
		isMutual, err = h.db.IsFollowing(ctx, req.FollowingId, req.FollowerId)
		if err != nil {
			return &follow.CheckFollowStatusResponse{
				Code:    constants.DatabaseErrorCode,
				Message: fmt.Sprintf("查询互关状态失败: %s", err),
			}, nil
		}
	}

	return &follow.CheckFollowStatusResponse{
		Code:        constants.SuccessCode,
		Message:     "查询成功",
		IsFollowing: isFollowing,
		IsMutual:    isMutual,
	}, nil
}
