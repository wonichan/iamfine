package handler

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"gorm.io/gorm"
	"hupu/kitex_gen/post"
	"hupu/shared/models"
)

type PostHandler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewPostHandler(db *gorm.DB, rdb *redis.Client) *PostHandler {
	return &PostHandler{
		db:  db,
		rdb: rdb,
	}
}

func (h *PostHandler) CreatePost(ctx context.Context, req *post.CreatePostRequest) (*post.CreatePostResponse, error) {
	// 生成帖子ID
	postID := xid.New().String()

	// 创建帖子
	newPost := models.Post{
		ID:      postID,
		UserID:  req.UserId,
		Title:   req.Title,
		Content: req.Content,
		Images:  models.StringArray(req.Images),
	}

	err := h.db.Create(&newPost).Error
	if err != nil {
		return &post.CreatePostResponse{
			Code:    500,
			Message: "创建帖子失败",
		}, err
	}

	return &post.CreatePostResponse{
		Code:    200,
		Message: "创建成功",
		Post: &post.Post{
			Id:           newPost.ID,
			UserId:       newPost.UserID,
			Title:        newPost.Title,
			Content:      newPost.Content,
			Images:       []string(newPost.Images),
			LikeCount:    int32(newPost.LikeCount),
			CommentCount: int32(newPost.CommentCount),
			CreatedAt:    newPost.CreatedAt.Unix(),
			UpdatedAt:    newPost.UpdatedAt.Unix(),
		},
	}, nil
}

func (h *PostHandler) GetPost(ctx context.Context, req *post.GetPostRequest) (*post.GetPostResponse, error) {
	var postModel models.Post
	err := h.db.Where("id = ?", req.PostId).First(&postModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &post.GetPostResponse{
				Code:    404,
				Message: "帖子不存在",
			}, nil
		}
		return &post.GetPostResponse{
			Code:    500,
			Message: "查询帖子失败",
		}, err
	}

	return &post.GetPostResponse{
		Code:    200,
		Message: "查询成功",
		Post: &post.Post{
			Id:           postModel.ID,
			UserId:       postModel.UserID,
			Title:        postModel.Title,
			Content:      postModel.Content,
			Images:       []string(postModel.Images),
			LikeCount:    int32(postModel.LikeCount),
			CommentCount: int32(postModel.CommentCount),
			CreatedAt:    postModel.CreatedAt.Unix(),
			UpdatedAt:    postModel.UpdatedAt.Unix(),
		},
	}, nil
}

func (h *PostHandler) GetPostList(ctx context.Context, req *post.GetPostListRequest) (*post.GetPostListResponse, error) {
	var posts []models.Post
	var total int64

	query := h.db.Model(&models.Post{})
	if req.UserId != nil {
		query = query.Where("user_id = ?", *req.UserId)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return &post.GetPostListResponse{
			Code:    500,
			Message: "查询失败",
		}, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err = query.Offset(int(offset)).Limit(int(req.PageSize)).Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return &post.GetPostListResponse{
			Code:    500,
			Message: "查询失败",
		}, err
	}

	// 转换数据格式
	var postList []*post.Post
	for _, p := range posts {
		postList = append(postList, &post.Post{
			Id:           p.ID,
			UserId:       p.UserID,
			Title:        p.Title,
			Content:      p.Content,
			Images:       []string(p.Images),
			LikeCount:    int32(p.LikeCount),
			CommentCount: int32(p.CommentCount),
			CreatedAt:    p.CreatedAt.Unix(),
			UpdatedAt:    p.UpdatedAt.Unix(),
		})
	}

	return &post.GetPostListResponse{
		Code:    200,
		Message: "查询成功",
		Posts:   postList,
		Total:   int32(total),
	}, nil
}
