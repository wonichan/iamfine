package handler

import (
	"context"
	"hupu/kitex_gen/post"
	"hupu/services/post/repository"
	"hupu/shared/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type PostHandler struct {
	db  repository.PostRepository
	rdb repository.PostRepository
}

func NewPostHandler(db *gorm.DB, rdb *redis.Client) *PostHandler {
	return &PostHandler{
		db:  repository.NewPostRepository(db),
		rdb: repository.NewPostRedisRepo(rdb),
	}
}

func (h *PostHandler) CreatePost(ctx context.Context, req *post.CreatePostRequest) (*post.CreatePostResponse, error) {
	// 创建帖子
	newPost := &models.Post{
		UserID:  req.UserId,
		Title:   req.Title,
		Content: req.Content,
		Images:  models.StringArray(req.Images),
	}

	err := h.rdb.CreatePost(ctx, newPost)
	if err != nil {
		return &post.CreatePostResponse{
			Code:    500,
			Message: "创建帖子失败",
		}, err
	}

	return &post.CreatePostResponse{
		Code:    0,
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
	postModel, err := h.rdb.GetPost(ctx, req.PostId)
	if err != nil {
		if err.Error() == "帖子不存在" {
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
	var posts []*models.Post
	var err error

	if req.UserId != nil {
		// 根据用户ID获取帖子列表
		posts, err = h.rdb.GetPostsByUserID(ctx, *req.UserId, req.Page, req.PageSize)
	} else {
		// 获取全部帖子列表
		posts, err = h.rdb.GetPostList(ctx, req.Page, req.PageSize)
	}

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
		Total:   int32(len(postList)),
	}, nil
}
