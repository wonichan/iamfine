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

// UpdatePost 更新帖子
func (h *PostHandler) UpdatePost(ctx context.Context, req *post.UpdatePostRequest) (*post.UpdatePostResponse, error) {
	// TODO: 实现更新帖子逻辑
	return &post.UpdatePostResponse{
		Code:    0,
		Message: "更新成功",
	}, nil
}

// DeletePost 删除帖子
func (h *PostHandler) DeletePost(ctx context.Context, req *post.DeletePostRequest) (*post.DeletePostResponse, error) {
	// TODO: 实现删除帖子逻辑
	return &post.DeletePostResponse{
		Code:    0,
		Message: "删除成功",
	}, nil
}

// GetRecommendPosts 获取推荐帖子
func (h *PostHandler) GetRecommendPosts(ctx context.Context, req *post.GetRecommendPostsRequest) (*post.GetRecommendPostsResponse, error) {
	// TODO: 实现获取推荐帖子逻辑
	return &post.GetRecommendPostsResponse{
		Code:    0,
		Message: "获取成功",
		Posts:   []*post.Post{},
		Total:   0,
		HasMore: false,
	}, nil
}

// GetHotPosts 获取热门帖子
func (h *PostHandler) GetHotPosts(ctx context.Context, req *post.GetHotPostsRequest) (*post.GetHotPostsResponse, error) {
	// TODO: 实现获取热门帖子逻辑
	return &post.GetHotPostsResponse{
		Code:    0,
		Message: "获取成功",
		Posts:   []*post.Post{},
		Total:   0,
		HasMore: false,
	}, nil
}

// GetHighScorePosts 获取高分帖子
func (h *PostHandler) GetHighScorePosts(ctx context.Context, req *post.GetHighScorePostsRequest) (*post.GetHighScorePostsResponse, error) {
	// TODO: 实现获取高分帖子逻辑
	return &post.GetHighScorePostsResponse{
		Code:    0,
		Message: "获取成功",
		Posts:   []*post.Post{},
		Total:   0,
		HasMore: false,
	}, nil
}

// GetLowScorePosts 获取低分帖子
func (h *PostHandler) GetLowScorePosts(ctx context.Context, req *post.GetLowScorePostsRequest) (*post.GetLowScorePostsResponse, error) {
	// TODO: 实现获取低分帖子逻辑
	return &post.GetLowScorePostsResponse{
		Code:    0,
		Message: "获取成功",
		Posts:   []*post.Post{},
		Total:   0,
		HasMore: false,
	}, nil
}

// GetControversialPosts 获取争议帖子
func (h *PostHandler) GetControversialPosts(ctx context.Context, req *post.GetControversialPostsRequest) (*post.GetControversialPostsResponse, error) {
	// TODO: 实现获取争议帖子逻辑
	return &post.GetControversialPostsResponse{
		Code:    0,
		Message: "获取成功",
		Posts:   []*post.Post{},
		Total:   0,
		HasMore: false,
	}, nil
}

// SearchPosts 搜索帖子
func (h *PostHandler) SearchPosts(ctx context.Context, req *post.SearchPostsRequest) (*post.SearchPostsResponse, error) {
	// TODO: 实现搜索帖子逻辑
	return &post.SearchPostsResponse{
		Code:    0,
		Message: "搜索成功",
		Posts:   []*post.Post{},
		Total:   0,
	}, nil
}

// GetTopic 获取话题详情
func (h *PostHandler) GetTopic(ctx context.Context, req *post.GetTopicRequest) (*post.GetTopicResponse, error) {
	// TODO: 实现获取话题详情逻辑
	return &post.GetTopicResponse{
		Code:    0,
		Message: "获取成功",
	}, nil
}

// GetHotTopics 获取热门话题
func (h *PostHandler) GetHotTopics(ctx context.Context, req *post.GetHotTopicsRequest) (*post.GetHotTopicsResponse, error) {
	// TODO: 实现获取热门话题逻辑
	return &post.GetHotTopicsResponse{
		Code:    0,
		Message: "获取成功",
		Topics:  []*post.Topic{},
	}, nil
}

// GetTopicCategories 获取话题分类
func (h *PostHandler) GetTopicCategories(ctx context.Context, req *post.GetTopicCategoriesRequest) (*post.GetTopicCategoriesResponse, error) {
	// TODO: 实现获取话题分类逻辑
	return &post.GetTopicCategoriesResponse{
		Code:    0,
		Message: "获取成功",
		Topics:  []*post.Topic{},
	}, nil
}

// SearchTopics 搜索话题
func (h *PostHandler) SearchTopics(ctx context.Context, req *post.SearchTopicsRequest) (*post.SearchTopicsResponse, error) {
	// TODO: 实现搜索话题逻辑
	return &post.SearchTopicsResponse{
		Code:    0,
		Message: "搜索成功",
		Topics:  []*post.Topic{},
		Total:   0,
	}, nil
}

// ShareTopic 分享话题
func (h *PostHandler) ShareTopic(ctx context.Context, req *post.ShareTopicRequest) (*post.ShareTopicResponse, error) {
	// TODO: 实现分享话题逻辑
	return &post.ShareTopicResponse{
		Code:    0,
		Message: "分享成功",
	}, nil
}

// GetUserRating 获取用户评分
func (h *PostHandler) GetUserRating(ctx context.Context, req *post.GetUserRatingRequest) (*post.GetUserRatingResponse, error) {
	// TODO: 实现获取用户评分逻辑
	return &post.GetUserRatingResponse{
		Code:     0,
		Message:  "获取成功",
		IsRated:  false,
	}, nil
}

// UpdateRating 更新评分
func (h *PostHandler) UpdateRating(ctx context.Context, req *post.UpdateRatingRequest) (*post.UpdateRatingResponse, error) {
	// TODO: 实现更新评分逻辑
	return &post.UpdateRatingResponse{
		Code:          0,
		Message:       "更新成功",
		AverageScore:  0,
		TotalRatings:  0,
	}, nil
}

// DeleteRating 删除评分
func (h *PostHandler) DeleteRating(ctx context.Context, req *post.DeleteRatingRequest) (*post.DeleteRatingResponse, error) {
	// TODO: 实现删除评分逻辑
	return &post.DeleteRatingResponse{
		Code:          0,
		Message:       "删除成功",
		AverageScore:  0,
		TotalRatings:  0,
	}, nil
}
