package handler

import (
	"context"
	"hupu/kitex_gen/like"
	"hupu/services/like/repository"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type LikeHandler struct {
	db  repository.LikeRepository
	rdb repository.LikeRepository
}

func NewLikeHandler(db *gorm.DB, rdb *redis.Client) *LikeHandler {
	return &LikeHandler{
		db:  repository.NewLikeRepository(db),
		rdb: repository.NewLikeRedisRepo(rdb),
	}
}

func (h *LikeHandler) Like(ctx context.Context, req *like.LikeRequest) (*like.LikeResponse, error) {
	err := h.rdb.Like(ctx, req.UserId, req.TargetId, req.TargetType)
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
	err := h.rdb.Unlike(ctx, req.UserId, req.TargetId, req.TargetType)
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
	isLiked, err := h.rdb.IsLiked(ctx, req.UserId, req.TargetId, req.TargetType)
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
	likes, err := h.rdb.GetLikeList(ctx, req.UserId, req.Page, req.PageSize)
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
