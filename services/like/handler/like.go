package handler

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"hupu/kitex_gen/like"
	"hupu/shared/models"
)

type LikeHandler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewLikeHandler(db *gorm.DB, rdb *redis.Client) *LikeHandler {
	return &LikeHandler{
		db:  db,
		rdb: rdb,
	}
}

func (h *LikeHandler) Like(ctx context.Context, req *like.LikeRequest) (*like.LikeResponse, error) {
	// 检查是否已经点赞
	var existLike models.Like
	err := h.db.Where("user_id = ? AND target_id = ? AND target_type = ?", req.UserId, req.TargetId, req.TargetType).First(&existLike).Error
	if err == nil {
		return &like.LikeResponse{
			Code:    400,
			Message: "已经点赞过了",
		}, nil
	}

	// 开始事务
	tx := h.db.Begin()

	// 创建点赞记录
	newLike := models.Like{
		UserID:     req.UserId,
		TargetID:   req.TargetId,
		TargetType: req.TargetType,
	}

	err = tx.Create(&newLike).Error
	if err != nil {
		tx.Rollback()
		return &like.LikeResponse{
			Code:    500,
			Message: "点赞失败",
		}, err
	}

	// 更新目标对象的点赞数
	if req.TargetType == "post" {
		err = tx.Model(&models.Post{}).Where("id = ?", req.TargetId).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
	} else if req.TargetType == "comment" {
		err = tx.Model(&models.Comment{}).Where("id = ?", req.TargetId).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
	}

	if err != nil {
		tx.Rollback()
		return &like.LikeResponse{
			Code:    500,
			Message: "更新点赞数失败",
		}, err
	}

	tx.Commit()

	return &like.LikeResponse{
		Code:    200,
		Message: "点赞成功",
	}, nil
}

func (h *LikeHandler) Unlike(ctx context.Context, req *like.UnlikeRequest) (*like.UnlikeResponse, error) {
	// 开始事务
	tx := h.db.Begin()

	// 删除点赞记录
	err := tx.Where("user_id = ? AND target_id = ? AND target_type = ?", req.UserId, req.TargetId, req.TargetType).Delete(&models.Like{}).Error
	if err != nil {
		tx.Rollback()
		return &like.UnlikeResponse{
			Code:    500,
			Message: "取消点赞失败",
		}, err
	}

	// 更新目标对象的点赞数
	if req.TargetType == "post" {
		err = tx.Model(&models.Post{}).Where("id = ?", req.TargetId).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
	} else if req.TargetType == "comment" {
		err = tx.Model(&models.Comment{}).Where("id = ?", req.TargetId).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
	}

	if err != nil {
		tx.Rollback()
		return &like.UnlikeResponse{
			Code:    500,
			Message: "更新点赞数失败",
		}, err
	}

	tx.Commit()

	return &like.UnlikeResponse{
		Code:    200,
		Message: "取消点赞成功",
	}, nil
}

func (h *LikeHandler) GetLikeList(ctx context.Context, req *like.GetLikeListRequest) (*like.GetLikeListResponse, error) {
	var likes []models.Like
	offset := (req.Page - 1) * req.PageSize
	err := h.db.Where("user_id = ?", req.UserId).Offset(int(offset)).Limit(int(req.PageSize)).Order("created_at DESC").Find(&likes).Error
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
		Code:    200,
		Message: "查询成功",
		Likes:   likeList,
	}, nil
}
