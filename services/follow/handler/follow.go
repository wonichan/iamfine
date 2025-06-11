package handler

import (
	"context"
	"hupu/kitex_gen/follow"
	"hupu/services/follow/repository"

	"github.com/go-redis/redis/v8"
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

	return &follow.GetFollowListResponse{
		Code:    0,
		Message: "查询成功",
		UserIds: followList,
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
