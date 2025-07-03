package handler

import (
	"context"
	"hupu/kitex_gen/follow"
	"hupu/services/follow/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type FollowHandler struct {
	db  repository.FollowRepository
	rdb repository.FollowRepository
}

func NewFollowHandler(db *gorm.DB, rdb *redis.Client) *FollowHandler {
	return &FollowHandler{
		db:  repository.NewFollowRepository(db),
		rdb: repository.NewFollowRedisRepo(rdb),
	}
}

func (h *FollowHandler) Follow(ctx context.Context, req *follow.FollowRequest) (*follow.FollowResponse, error) {
	err := h.rdb.Follow(ctx, req.FollowerId, req.FollowingId)
	if err != nil {
		return &follow.FollowResponse{
			Code:    400,
			Message: err.Error(),
		}, nil
	}

	return &follow.FollowResponse{
		Code:    0,
		Message: "关注成功",
	}, nil
}

func (h *FollowHandler) Unfollow(ctx context.Context, req *follow.UnfollowRequest) (*follow.UnfollowResponse, error) {
	err := h.rdb.Unfollow(ctx, req.FollowerId, req.FollowingId)
	if err != nil {
		return &follow.UnfollowResponse{
			Code:    500,
			Message: "取消关注失败",
		}, err
	}

	return &follow.UnfollowResponse{
		Code:    0,
		Message: "取消关注成功",
	}, nil
}

func (h *FollowHandler) IsFollowing(ctx context.Context, req *follow.FollowRequest) (*follow.FollowResponse, error) {
	isFollowing, err := h.rdb.IsFollowing(ctx, req.FollowerId, req.FollowingId)
	if err != nil {
		return &follow.FollowResponse{
			Code:    500,
			Message: "查询失败",
		}, err
	}

	if !isFollowing {
		return &follow.FollowResponse{
			Code:    400,
			Message: "未关注",
		}, nil
	}

	return &follow.FollowResponse{
		Code:        0,
		Message:     "已关注",
		IsFollowing: true,
	}, nil
}

func (h *FollowHandler) GetFollowList(ctx context.Context, req *follow.GetFollowListRequest) (*follow.GetFollowListResponse, error) {
	followList, err := h.rdb.GetFollowList(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &follow.GetFollowListResponse{
			Code:    500,
			Message: "查询关注列表失败",
		}, err
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
		Code:    0,
		Message: "查询成功",
		Follows: follows,
		Total:   int32(len(follows)),
		HasMore: len(follows) == int(req.PageSize),
	}, nil
}

func (h *FollowHandler) GetFollowerList(ctx context.Context, req *follow.GetFollowerListRequest) (*follow.GetFollowerListResponse, error) {
	follows, err := h.rdb.GetFollowerList(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &follow.GetFollowerListResponse{
			Code:    500,
			Message: "查询粉丝列表失败",
		}, err
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
		Code:      0,
		Message:   "查询成功",
		Followers: followerList,
	}, nil
}

// GetFollowCount 获取关注数量
func (h *FollowHandler) GetFollowCount(ctx context.Context, req *follow.GetFollowCountRequest) (*follow.GetFollowCountResponse, error) {
	count, err := h.rdb.GetFollowCount(ctx, req.UserId)
	if err != nil {
		return &follow.GetFollowCountResponse{
			Code:    500,
			Message: "查询关注数量失败",
		}, err
	}

	return &follow.GetFollowCountResponse{
		Code:    0,
		Message: "查询成功",
		Count:   int32(count),
	}, nil
}

// GetFollowerCount 获取粉丝数量
func (h *FollowHandler) GetFollowerCount(ctx context.Context, req *follow.GetFollowerCountRequest) (*follow.GetFollowerCountResponse, error) {
	count, err := h.rdb.GetFollowerCount(ctx, req.UserId)
	if err != nil {
		return &follow.GetFollowerCountResponse{
			Code:    500,
			Message: "查询粉丝数量失败",
		}, err
	}

	return &follow.GetFollowerCountResponse{
		Code:    0,
		Message: "查询成功",
		Count:   int32(count),
	}, nil
}

// GetMutualFollows 获取共同关注
func (h *FollowHandler) GetMutualFollows(ctx context.Context, req *follow.GetMutualFollowsRequest) (*follow.GetMutualFollowsResponse, error) {
	userIds, err := h.rdb.GetMutualFollows(ctx, req.UserId, req.TargetUserId, req.Page, req.PageSize)
	if err != nil {
		return &follow.GetMutualFollowsResponse{
			Code:    500,
			Message: "查询共同关注失败",
		}, err
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
		Code:    0,
		Message: "查询成功",
		MutualFollows: mutualFollows,
		Total:   int32(len(mutualFollows)),
		HasMore: len(mutualFollows) == int(req.PageSize),
	}, nil
}

// CheckFollowStatus 检查关注状态
func (h *FollowHandler) CheckFollowStatus(ctx context.Context, req *follow.CheckFollowStatusRequest) (*follow.CheckFollowStatusResponse, error) {
	// 检查是否关注
	isFollowing, err := h.rdb.IsFollowing(ctx, req.FollowerId, req.FollowingId)
	if err != nil {
		return &follow.CheckFollowStatusResponse{
			Code:    500,
			Message: "查询关注状态失败",
		}, err
	}

	// 如果已关注，检查是否互相关注
	var isMutual bool
	if isFollowing {
		isMutual, err = h.rdb.IsFollowing(ctx, req.FollowingId, req.FollowerId)
		if err != nil {
			return &follow.CheckFollowStatusResponse{
				Code:    500,
				Message: "查询互关状态失败",
			}, err
		}
	}

	return &follow.CheckFollowStatusResponse{
		Code:        0,
		Message:     "查询成功",
		IsFollowing: isFollowing,
		IsMutual:    isMutual,
	}, nil
}
