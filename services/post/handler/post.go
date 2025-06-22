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
		UserID:      req.UserId,
		Title:       req.Title,
		Content:     req.Content,
		TopicID:     req.TopicId,
		Tags:        req.Tags,
		IsAnonymous: req.IsAnonymous,
		Location:    req.Location,
		Category:    models.PostCategory(req.Category),
		Images:      models.StringArray(req.Images),
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
		Post:    models.PostToKitexPost(newPost),
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
		Post:    models.PostToKitexPost(postModel),
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
		postList = append(postList, models.PostToKitexPost(p))
	}

	return &post.GetPostListResponse{
		Code:    200,
		Message: "查询成功",
		Posts:   postList,
		Total:   int32(len(postList)),
	}, nil
}

// 话题管理相关方法
func (h *PostHandler) CreateTopic(ctx context.Context, req *post.CreateTopicRequest) (*post.CreateTopicResponse, error) {
	topic := &models.Topic{
		Name:        req.Name,
		Description: &req.Description,
	}

	err := h.rdb.CreateTopic(ctx, topic)
	if err != nil {
		return &post.CreateTopicResponse{
			Code:    500,
			Message: "创建话题失败",
		}, err
	}

	description := ""
	if topic.Description != nil {
		description = *topic.Description
	}

	return &post.CreateTopicResponse{
		Code:    0,
		Message: "创建成功",
		Topic: &post.Topic{
			Id:               topic.ID,
			Name:             topic.Name,
			Description:      description,
			Icon:             "",              // Topic模型中没有Icon字段
			Color:            "",              // Topic模型中没有Color字段
			ParticipantCount: topic.PostCount, // 使用PostCount作为参与者数量
			CreatedAt:        topic.CreatedAt.Unix(),
			UpdatedAt:        topic.UpdatedAt.Unix(),
		},
	}, nil
}

func (h *PostHandler) GetTopicList(ctx context.Context, req *post.GetTopicListRequest) (*post.GetTopicListResponse, error) {
	topics, err := h.rdb.GetTopicList(ctx, req.Page, req.PageSize)
	if err != nil {
		return &post.GetTopicListResponse{
			Code:    500,
			Message: "查询话题列表失败",
		}, err
	}

	// 转换数据格式
	var topicList []*post.Topic
	for _, t := range topics {
		description := ""
		if t.Description != nil {
			description = *t.Description
		}
		topicList = append(topicList, &post.Topic{
			Id:               t.ID,
			Name:             t.Name,
			Description:      description,
			Icon:             "",          // Topic模型中没有Icon字段
			Color:            "",          // Topic模型中没有Color字段
			ParticipantCount: t.PostCount, // 使用PostCount作为参与者数量
			CreatedAt:        t.CreatedAt.Unix(),
			UpdatedAt:        t.UpdatedAt.Unix(),
		})
	}

	return &post.GetTopicListResponse{
		Code:    0,
		Message: "查询成功",
		Topics:  topicList,
		Total:   int32(len(topicList)),
	}, nil
}

// 收藏功能相关方法
func (h *PostHandler) CollectPost(ctx context.Context, req *post.CollectPostRequest) (*post.CollectPostResponse, error) {
	err := h.rdb.FavoritePost(ctx, req.UserId, req.PostId)
	if err != nil {
		return &post.CollectPostResponse{
			Code:    500,
			Message: "收藏失败",
		}, err
	}

	return &post.CollectPostResponse{
		Code:    0,
		Message: "收藏成功",
	}, nil
}

func (h *PostHandler) UncollectPost(ctx context.Context, req *post.UncollectPostRequest) (*post.UncollectPostResponse, error) {
	err := h.rdb.UnfavoritePost(ctx, req.UserId, req.PostId)
	if err != nil {
		return &post.UncollectPostResponse{
			Code:    500,
			Message: "取消收藏失败",
		}, err
	}

	return &post.UncollectPostResponse{
		Code:    0,
		Message: "取消收藏成功",
	}, nil
}

func (h *PostHandler) GetCollectedPosts(ctx context.Context, req *post.GetCollectedPostsRequest) (*post.GetCollectedPostsResponse, error) {
	posts, err := h.rdb.GetFavoriteList(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &post.GetCollectedPostsResponse{
			Code:    500,
			Message: "查询收藏列表失败",
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

	return &post.GetCollectedPostsResponse{
		Code:    0,
		Message: "查询成功",
		Posts:   postList,
		Total:   int32(len(postList)),
	}, nil
}

// 评分功能相关方法
func (h *PostHandler) RatePost(ctx context.Context, req *post.RatePostRequest) (*post.RatePostResponse, error) {
	rating := &models.PostRating{
		UserID:  req.UserId,
		PostID:  req.PostId,
		Score:   int32(req.Score), // 将float64转换为int32
		Comment: req.Comment,      // req.Comment已经是*string类型
	}

	err := h.rdb.RatePost(ctx, rating)
	if err != nil {
		return &post.RatePostResponse{
			Code:    500,
			Message: "评分失败",
		}, err
	}

	// TODO: 计算平均分和总评分数
	// 这里需要从数据库查询相关统计信息
	return &post.RatePostResponse{
		Code:         0,
		Message:      "评分成功",
		AverageScore: req.Score, // 临时返回当前评分
		TotalRatings: 1,         // 临时返回1
	}, nil
}

func (h *PostHandler) GetRatingRank(ctx context.Context, req *post.GetRatingRankRequest) (*post.GetRatingRankResponse, error) {
	posts, err := h.rdb.GetScoreRanking(ctx, req.Page, req.PageSize)
	if err != nil {
		return &post.GetRatingRankResponse{
			Code:    500,
			Message: "查询评分排行榜失败",
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
			// Score字段在Post模型中不存在，这里设置为0或从其他地方计算
			Score:     0.0,
			CreatedAt: p.CreatedAt.Unix(),
			UpdatedAt: p.UpdatedAt.Unix(),
		})
	}

	return &post.GetRatingRankResponse{
		Code:    0,
		Message: "查询成功",
		Posts:   postList,
		Total:   int32(len(postList)),
	}, nil
}
