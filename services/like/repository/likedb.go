package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"hupu/shared/models"
	"hupu/shared/utils"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository() *LikeRepository {
	return &LikeRepository{
		db: utils.GetDB(),
	}
}

func (lr *LikeRepository) Like(ctx context.Context, userID, targetID, targetType string) error {
	// 检查是否已经点赞
	var existLike models.Like
	err := lr.db.Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).First(&existLike).Error
	if err == nil {
		return fmt.Errorf("已经点赞过了")
	}

	// 开始事务
	tx := lr.db.Begin()

	// 创建点赞记录
	newLike := models.Like{
		UserID:     userID,
		TargetID:   targetID,
		TargetType: targetType,
	}

	err = tx.Create(&newLike).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 更新目标对象的点赞数
	if targetType == "post" {
		err = tx.Model(&models.Post{}).Where("id = ?", targetID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
	} else if targetType == "comment" {
		err = tx.Model(&models.Comment{}).Where("id = ?", targetID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (lr *LikeRepository) Unlike(ctx context.Context, userID, targetID, targetType string) error {
	// 开始事务
	tx := lr.db.Begin()

	// 删除点赞记录
	err := tx.Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).Delete(&models.Like{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 更新目标对象的点赞数
	if targetType == "post" {
		err = tx.Model(&models.Post{}).Where("id = ?", targetID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
	} else if targetType == "comment" {
		err = tx.Model(&models.Comment{}).Where("id = ?", targetID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (lr *LikeRepository) IsLiked(ctx context.Context, userID, targetID, targetType string) (bool, error) {
	var existLike models.Like
	err := lr.db.Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).First(&existLike).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (lr *LikeRepository) GetLikeList(ctx context.Context, userID string, page, pageSize int32) ([]*models.Like, error) {
	var likes []models.Like
	offset := (page - 1) * pageSize
	err := lr.db.Where("user_id = ?", userID).Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&likes).Error
	if err != nil {
		return nil, err
	}

	// 转换为指针切片
	var likeList []*models.Like
	for i := range likes {
		likeList = append(likeList, &likes[i])
	}

	return likeList, nil
}

func (lr *LikeRepository) GetLikeCount(ctx context.Context, targetID, targetType string) (int64, error) {
	var count int64
	err := lr.db.Model(&models.Like{}).Where("target_id = ? AND target_type = ?", targetID, targetType).Count(&count).Error
	return count, err
}

func (lr *LikeRepository) GetLikeUsers(ctx context.Context, targetID, targetType string, page, pageSize int32) ([]string, error) {
	var userIDs []string
	offset := (page - 1) * pageSize
	err := lr.db.Model(&models.Like{}).
		Select("user_id").
		Where("target_id = ? AND target_type = ?", targetID, targetType).
		Offset(int(offset)).
		Limit(int(pageSize)).
		Order("created_at DESC").
		Pluck("user_id", &userIDs).Error

	return userIDs, err
}
