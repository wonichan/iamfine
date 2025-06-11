package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"

	"hupu/shared/models"
)

type commentRedisRepo struct {
	rdb *redis.Client
}

func NewCommentRedisRepo(rdb *redis.Client) CommentRepository {
	return &commentRedisRepo{
		rdb: rdb,
	}
}

func (cr *commentRedisRepo) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	// 生成评论ID
	commentID := xid.New().String()

	// 创建评论
	newComment := &models.Comment{
		ID:        commentID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		ParentID:  comment.ParentID,
		LikeCount: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存到Redis
	err := cr.saveComment(ctx, newComment)
	if err != nil {
		return nil, err
	}

	return newComment, nil
}

func (cr *commentRedisRepo) GetCommentList(ctx context.Context, postID string, parentID *string, page, pageSize int32) ([]*models.Comment, error) {
	// 构建查询键
	var listKey string
	if parentID != nil {
		listKey = fmt.Sprintf("comment:list:post:%s:parent:%s", postID, *parentID)
	} else {
		listKey = fmt.Sprintf("comment:list:post:%s:top", postID)
	}

	// 分页获取评论ID列表
	start := (page - 1) * pageSize
	end := start + pageSize - 1
	commentIDs, err := cr.rdb.LRange(ctx, listKey, int64(start), int64(end)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []*models.Comment{}, nil
		}
		return nil, err
	}

	// 批量获取评论详情
	var comments []*models.Comment
	for _, commentID := range commentIDs {
		comment, err := cr.GetComment(ctx, commentID)
		if err == nil && comment != nil {
			comments = append(comments, comment)
		}
	}

	return comments, nil
}

func (cr *commentRedisRepo) GetComment(ctx context.Context, commentID string) (*models.Comment, error) {
	commentKey := fmt.Sprintf("comment:%s", commentID)
	commentData, err := cr.rdb.Get(ctx, commentKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("评论不存在")
		}
		return nil, err
	}

	// 反序列化评论数据
	var comment models.Comment
	err = json.Unmarshal([]byte(commentData), &comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (cr *commentRedisRepo) UpdateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	// 检查评论是否存在
	existComment, err := cr.GetComment(ctx, comment.ID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if comment.Content != "" {
		existComment.Content = comment.Content
	}
	existComment.UpdatedAt = time.Now()

	// 保存更新后的评论
	err = cr.saveComment(ctx, existComment)
	if err != nil {
		return nil, err
	}

	return existComment, nil
}

func (cr *commentRedisRepo) DeleteComment(ctx context.Context, commentID string) error {
	// 获取评论信息
	comment, err := cr.GetComment(ctx, commentID)
	if err != nil {
		return err
	}

	// 删除评论数据
	commentKey := fmt.Sprintf("comment:%s", commentID)

	// 构建列表键
	var listKey string
	if comment.ParentID != nil && *comment.ParentID != "" {
		listKey = fmt.Sprintf("comment:list:post:%s:parent:%s", comment.PostID, *comment.ParentID)
	} else {
		listKey = fmt.Sprintf("comment:list:post:%s:top", comment.PostID)
	}

	pipe := cr.rdb.Pipeline()
	pipe.Del(ctx, commentKey)
	pipe.LRem(ctx, listKey, 0, commentID)
	_, err = pipe.Exec(ctx)

	return err
}

// 辅助方法：保存评论到Redis
func (cr *commentRedisRepo) saveComment(ctx context.Context, comment *models.Comment) error {
	// 序列化评论数据
	commentData, err := json.Marshal(comment)
	if err != nil {
		return err
	}

	// 构建键名
	commentKey := fmt.Sprintf("comment:%s", comment.ID)

	// 构建列表键
	var listKey string
	if comment.ParentID != nil && *comment.ParentID != "" {
		listKey = fmt.Sprintf("comment:list:post:%s:parent:%s", comment.PostID, *comment.ParentID)
	} else {
		listKey = fmt.Sprintf("comment:list:post:%s:top", comment.PostID)
	}

	// 使用Pipeline批量操作
	pipe := cr.rdb.Pipeline()
	pipe.Set(ctx, commentKey, commentData, 24*time.Hour)
	pipe.LPush(ctx, listKey, comment.ID)
	pipe.Expire(ctx, listKey, 24*time.Hour)
	_, err = pipe.Exec(ctx)

	return err
}