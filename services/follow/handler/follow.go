package handler

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"hupu/kitex_gen/follow"
	"hupu/shared/models"
)

type FollowHandler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewFollowHandler(db *gorm.DB, rdb *redis.Client) *FollowHandler {
	return &FollowHandler{
		db:  db,
		rdb: rdb,
	}
}

func (h *FollowHandler) Follow(ctx context.Context, req *follow.FollowRequest) (*follow.FollowResponse, error) {
	// 检查是否已经关注
	var existFollow models.Follow
	err := h.db.Where("follower_id = ? AND following_id = ?", req.FollowerId, req.FollowingId).First(&existFollow).Error
	if err == nil {
		return &follow.FollowResponse{
			Code:    400,
			Message: "已经关注过了",
		}, nil
	}

	// 不能关注自己
	if req.FollowerId == req.FollowingId {
		return &follow.FollowResponse{
			Code:    400,
			Message: "不能关注自己",
		}, nil
	}

	// 创建关注关系
	newFollow := models.Follow{
		FollowerID:  req.FollowerId,
		FollowingID: req.FollowingId,
	}

	err = h.db.Create(&newFollow).Error
	if err != nil {
		return &follow.FollowResponse{
			Code:    500,
			Message: "关注失败",
		}, err
	}

	return &follow.FollowResponse{
		Code:    200,
		Message: "关注成功",
	}, nil
}

func (h *FollowHandler) Unfollow(ctx context.Context, req *follow.UnfollowRequest) (*follow.UnfollowResponse, error) {
	// 删除关注关系
	err := h.db.Where("follower_id = ? AND following_id = ?", req.FollowerId, req.FollowingId).Delete(&models.Follow{}).Error
	if err != nil {
		return &follow.UnfollowResponse{
			Code:    500,
			Message: "取消关注失败",
		}, err
	}

	return &follow.UnfollowResponse{
		Code:    200,
		Message: "取消关注成功",
	}, nil
}

func (h *FollowHandler) GetFollowList(ctx context.Context, req *follow.GetFollowListRequest) (*follow.GetFollowListResponse, error) {
	var follows []models.Follow
	offset := (req.Page - 1) * req.PageSize
	err := h.db.Where("follower_id = ?", req.UserId).Offset(int(offset)).Limit(int(req.PageSize)).Find(&follows).Error
	if err != nil {
		return &follow.GetFollowListResponse{
			Code:    500,
			Message: "查询关注列表失败",
		}, err
	}

	// 转换为响应格式
	var followList []*follow.Follow
	for _, f := range follows {
		followList = append(followList, &follow.Follow{
			FollowerId:  f.FollowerID,
			FollowingId: f.FollowingID,
			CreatedAt:   f.CreatedAt.Unix(),
		})
	}

	return &follow.GetFollowListResponse{
		Code:    200,
		Message: "查询成功",
		Follows: followList,
	}, nil
}

func (h *FollowHandler) GetFollowerList(ctx context.Context, req *follow.GetFollowerListRequest) (*follow.GetFollowerListResponse, error) {
	var follows []models.Follow
	offset := (req.Page - 1) * req.PageSize
	err := h.db.Where("following_id = ?", req.UserId).Offset(int(offset)).Limit(int(req.PageSize)).Find(&follows).Error
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
		Code:      200,
		Message:   "查询成功",
		Followers: followerList,
	}, nil
}
