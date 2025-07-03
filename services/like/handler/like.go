package handler

import (
	"context"
	"hupu/kitex_gen/like"
	"hupu/services/like/repository"

	"gorm.io/gorm"
)

type LikeHandler struct {
	db repository.LikeRepository
}

func NewLikeHandler(db *gorm.DB) *LikeHandler {
	return &LikeHandler{
		db: repository.NewLikeRepository(db),
	}
}

func (h *LikeHandler) Like(ctx context.Context, req *like.LikeRequest) (*like.LikeResponse, error) {
	err := h.db.Like(ctx, req.UserId, req.TargetId, req.TargetType)
	if err != nil {
		return &like.LikeResponse{
			Code:    400,
			Message: err.Error(),
		}, nil
	}

	return &like.LikeResponse{
		Code:    0,
		Message: "点赞成功",
	}, nil
}

func (h *LikeHandler) Unlike(ctx context.Context, req *like.UnlikeRequest) (*like.UnlikeResponse, error) {
	err := h.db.Unlike(ctx, req.UserId, req.TargetId, req.TargetType)
	if err != nil {
		return &like.UnlikeResponse{
			Code:    500,
			Message: "取消点赞失败",
		}, err
	}

	return &like.UnlikeResponse{
		Code:    0,
		Message: "取消点赞成功",
	}, nil
}

func (h *LikeHandler) IsLiked(ctx context.Context, req *like.LikeRequest) (*like.LikeResponse, error) {
	isLiked, err := h.db.IsLiked(ctx, req.UserId, req.TargetId, req.TargetType)
	if err != nil {
		return &like.LikeResponse{
			Code:    500,
			Message: "查询失败",
		}, err
	}

	if !isLiked {
		return &like.LikeResponse{
			Code:    400,
			Message: "未点赞",
		}, nil
	}

	return &like.LikeResponse{
		Code:    0,
		Message: "已点赞",
	}, nil
}

func (h *LikeHandler) GetLikeList(ctx context.Context, req *like.GetLikeListRequest) (*like.GetLikeListResponse, error) {
	likes, err := h.db.GetLikeList(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &like.GetLikeListResponse{
			Code:    500,
			Message: "查询点赞列表失败",
		}, err
	}

	// 转换为响应格式
	var likeList []*like.Like
	for _, l := range likes {
		likeList = append(likeList, &like.Like{
			UserId:     l.UserID,
			TargetId:   l.TargetID,
			TargetType: l.TargetType,
			CreatedAt:  l.CreatedAt.Unix(),
		})
	}

	return &like.GetLikeListResponse{
		Code:    0,
		Message: "查询成功",
		Likes:   likeList,
	}, nil
}

// GetLikeCount 获取点赞数量
func (h *LikeHandler) GetLikeCount(ctx context.Context, req *like.GetLikeCountRequest) (*like.GetLikeCountResponse, error) {
	count, err := h.db.GetLikeCount(ctx, req.TargetId, req.TargetType)
	if err != nil {
		return &like.GetLikeCountResponse{
			Code:    500,
			Message: "查询点赞数量失败",
		}, err
	}

	return &like.GetLikeCountResponse{
		Code:    0,
		Message: "查询成功",
		Count:   count,
	}, nil
}

// GetLikeUsers 获取点赞用户列表
func (h *LikeHandler) GetLikeUsers(ctx context.Context, req *like.GetLikeUsersRequest) (*like.GetLikeUsersResponse, error) {
	userIDs, err := h.db.GetLikeUsers(ctx, req.TargetId, req.TargetType, req.Page, req.PageSize)
	if err != nil {
		return &like.GetLikeUsersResponse{
			Code:    500,
			Message: "查询点赞用户失败",
		}, err
	}

	// 转换为响应格式
	var userList []*like.LikeUser
	for _, userID := range userIDs {
		// TODO: 这里应该调用用户服务获取用户详细信息
		// 目前先返回基本信息
		userList = append(userList, &like.LikeUser{
			UserId:    userID,
			Username:  "", // 需要从用户服务获取
			Nickname:  "", // 需要从用户服务获取
			Avatar:    "", // 需要从用户服务获取
			CreatedAt: 0,  // 需要从点赞记录获取
		})
	}

	return &like.GetLikeUsersResponse{
		Code:    0,
		Message: "查询成功",
		Users:   userList,
	}, nil
}
