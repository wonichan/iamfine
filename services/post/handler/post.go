package handler

import (
	"context"
	"fmt"
	"hupu/kitex_gen/post"
	"hupu/services/post/repository"
	"hupu/shared/constants"
	"hupu/shared/log"
	"hupu/shared/middleware"
	"hupu/shared/models"
	"hupu/shared/utils"
	"strings"

	"gorm.io/gorm"
)

type PostHandler struct {
	db *repository.PostRepository
}

func NewPostHandler() *PostHandler {
	return &PostHandler{
		db: repository.NewPostRepository(),
	}
}

func (h *PostHandler) CreatePost(ctx context.Context, req *post.CreatePostRequest) (*post.CreatePostResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Infof("CreatePost req: %v", req)

	// 参数验证
	validationErrors := middleware.ValidatePostCreation(req.UserId, req.Title, req.Content)
	if len(validationErrors) > 0 {
		errorMsg := "参数验证失败: "
		for _, ve := range validationErrors {
			errorMsg += fmt.Sprintf("%s: %s; ", ve.Field, ve.Message)
		}
		return &post.CreatePostResponse{
			Code:    constants.ValidationErrorCode,
			Message: errorMsg,
		}, nil
	}

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

	err := h.db.CreatePost(ctx, newPost)
	if err != nil {
		logger.Errorf("CreatePost failed: %s", err)
		return &post.CreatePostResponse{
			Code:    constants.PostCreateFailCode,
			Message: fmt.Sprintf("failed to create post: %s", err),
		}, nil
	}

	return &post.CreatePostResponse{
		Code: constants.SuccessCode,
		Post: models.PostToKitexPost(newPost),
	}, nil
}

func (h *PostHandler) GetPost(ctx context.Context, req *post.GetPostRequest) (*post.GetPostResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Infof("GetPost req: %v", req)

	// 参数验证
	if req.PostId == "" {
		return &post.GetPostResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to get post: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	postModel, err := h.db.GetPost(ctx, req.PostId)
	if err != nil {
		logger.Errorf("GetPost failed: %s", err)
		return &post.GetPostResponse{
			Code:    constants.PostNotFoundCode,
			Message: fmt.Sprintf("failed to get post: %s", err),
		}, nil
	}

	return &post.GetPostResponse{
		Code: constants.SuccessCode,
		Post: models.PostToKitexPost(postModel),
	}, nil
}

func (h *PostHandler) GetPostList(ctx context.Context, req *post.GetPostListRequest) (*post.GetPostListResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Infof("GetPostList req: %v", req)
	var posts []*models.Post
	var err error

	// 构建查询条件
	conditions := make(map[string]interface{})
	useConditions := false

	// 添加用户ID条件
	if req.UserId != nil && *req.UserId != "" {
		conditions[constants.ParamUserID] = *req.UserId
		useConditions = true
	}

	// 添加话题ID条件
	if req.TopicId != nil && *req.TopicId != "" {
		conditions[constants.ParamTopicID] = *req.TopicId
		useConditions = true
	}

	// 添加分类条件
	if req.Category != nil && int32(*req.Category) > 0 {
		conditions[constants.ParamCategory] = int32(*req.Category)
		useConditions = true
	}

	// 获取排序类型
	sortType := constants.SortTypeLatest // 默认按最新排序
	if req.SortType != nil && *req.SortType != "" {
		sortType = *req.SortType
		useConditions = true
	}

	// 根据是否有条件选择查询方法
	if useConditions {
		// 使用条件查询
		posts, err = h.db.GetPostListWithConditions(ctx, conditions, req.Page, req.PageSize, sortType)
	} else {
		// 无条件查询，先尝试缓存
		// posts, err = h.rdb.GetPostList(ctx, req.Page, req.PageSize)
		// 缓存失败或无数据，回退到数据库
		logger.Warnf("Cache failed, fallback to database: %s", err)
		posts, err = h.db.GetPostList(ctx, req.Page, req.PageSize)
	}

	if err != nil {
		logger.Errorf("GetPostList failed: %s", err)
		return &post.GetPostListResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("failed to get post list: %s", err),
		}, nil
	}

	// 转换数据格式
	var postList []*post.Post
	for _, p := range posts {
		postList = append(postList, models.PostToKitexPost(p))
	}

	return &post.GetPostListResponse{
		Code:  constants.SuccessCode,
		Posts: postList,
		Total: int32(len(postList)),
	}, nil
}

// 话题管理相关方法
func (h *PostHandler) CreateTopic(ctx context.Context, req *post.CreateTopicRequest) (*post.CreateTopicResponse, error) {
	// 参数验证
	if req.Name == "" {
		return &post.CreateTopicResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to create topic: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	topic := &models.Topic{
		Name:        req.Name,
		Description: &req.Description,
	}

	err := h.db.CreateTopic(ctx, topic)
	if err != nil {
		// 根据错误类型返回不同的错误码
		if strings.Contains(err.Error(), "already exists") {
			return &post.CreateTopicResponse{
				Code:    constants.TopicAlreadyExistsCode,
				Message: fmt.Sprintf("failed to create topic: %s", err),
			}, nil
		}
		return &post.CreateTopicResponse{
			Code:    constants.TopicCreateFailCode,
			Message: fmt.Sprintf("failed to create topic: %s", err),
		}, nil
	}

	description := ""
	if topic.Description != nil {
		description = *topic.Description
	}

	return &post.CreateTopicResponse{
		Code: constants.SuccessCode,
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
	topics, err := h.db.GetTopicList(ctx, req.Page, req.PageSize)
	if err != nil {
		return &post.GetTopicListResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("failed to get topic list: %s", err),
		}, nil
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
		Code:   constants.SuccessCode,
		Topics: topicList,
		Total:  int32(len(topicList)),
	}, nil
}

// 收藏功能相关方法
func (h *PostHandler) CollectPost(ctx context.Context, req *post.CollectPostRequest) (*post.CollectPostResponse, error) {
	// 参数验证
	if req.UserId == "" || req.PostId == "" {
		return &post.CollectPostResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to collect post: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	err := h.db.FavoritePost(ctx, req.UserId, req.PostId)
	if err != nil {
		// 根据错误类型返回不同的错误码
		if strings.Contains(err.Error(), "post not found") {
			return &post.CollectPostResponse{
				Code:    constants.PostNotFoundCode,
				Message: fmt.Sprintf("failed to collect post: %s", err),
			}, nil
		}
		if strings.Contains(err.Error(), "already favorited") {
			return &post.CollectPostResponse{
				Code:    constants.PostAlreadyFavoritedCode,
				Message: fmt.Sprintf("failed to collect post: %s", err),
			}, nil
		}
		return &post.CollectPostResponse{
			Code:    constants.PostFavoriteFailCode,
			Message: fmt.Sprintf("failed to collect post: %s", err),
		}, nil
	}

	return &post.CollectPostResponse{
		Code:    constants.SuccessCode,
		Message: "收藏成功",
	}, nil
}

func (h *PostHandler) UncollectPost(ctx context.Context, req *post.UncollectPostRequest) (*post.UncollectPostResponse, error) {
	// 参数验证
	if req.UserId == "" || req.PostId == "" {
		return &post.UncollectPostResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to uncollect post: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	err := h.db.UnfavoritePost(ctx, req.UserId, req.PostId)
	if err != nil {
		return &post.UncollectPostResponse{
			Code:    constants.PostUnfavoriteFailCode,
			Message: fmt.Sprintf("failed to uncollect post: %s", err),
		}, nil
	}

	return &post.UncollectPostResponse{
		Code: constants.SuccessCode,
	}, nil
}

func (h *PostHandler) GetCollectedPosts(ctx context.Context, req *post.GetCollectedPostsRequest) (*post.GetCollectedPostsResponse, error) {
	// 参数验证
	if req.UserId == "" {
		return &post.GetCollectedPostsResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to get collected posts: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	posts, err := h.db.GetFavoriteList(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		return &post.GetCollectedPostsResponse{
			Code:    constants.DatabaseErrorCode,
			Message: fmt.Sprintf("failed to get collected posts: %s", err),
		}, nil
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
		Code:  constants.SuccessCode,
		Posts: postList,
		Total: int32(len(postList)),
	}, nil
}

// 评分功能相关方法
func (h *PostHandler) RatePost(ctx context.Context, req *post.RatePostRequest) (*post.RatePostResponse, error) {
	// 参数验证
	if req.UserId == "" || req.PostId == "" || req.Score < 1 || req.Score > 5 {
		return &post.RatePostResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to rate post: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	// 检查帖子是否存在并获取帖子信息
	postInfo, err := h.db.GetPost(ctx, req.PostId)
	if err != nil {
		return &post.RatePostResponse{
			Code:    constants.PostNotFoundCode,
			Message: fmt.Sprintf("failed to rate post: %s", constants.GetErrorMessage(constants.PostNotFoundCode)),
		}, nil
	}

	// 防止用户给自己的帖子评分
	if postInfo.UserID == req.UserId {
		return &post.RatePostResponse{
			Code:    constants.ValidationErrorCode,
			Message: "不能给自己的帖子评分",
		}, nil
	}

	rating := &models.PostRating{
		UserID:  req.UserId,
		PostID:  req.PostId,
		Score:   int32(req.Score), // 将float64转换为int32
		Comment: req.Comment,      // req.Comment已经是*string类型
	}

	err = h.db.RatePost(ctx, rating)
	if err != nil {
		return &post.RatePostResponse{
			Code:    constants.PostRatingFailCode,
			Message: fmt.Sprintf("failed to rate post: %s", err),
		}, nil
	}

	// 获取帖子的评分统计信息
	averageScore, totalRatings, err := h.db.GetPostRatingStats(ctx, req.PostId)
	if err != nil {
		return &post.RatePostResponse{
			Code:    constants.PostRatingFailCode,
			Message: fmt.Sprintf("failed to get rating stats: %s", err),
		}, nil
	}

	return &post.RatePostResponse{
		Code:         constants.SuccessCode,
		AverageScore: averageScore,
		TotalRatings: totalRatings,
	}, nil
}

func (h *PostHandler) GetRatingRank(ctx context.Context, req *post.GetRatingRankRequest) (*post.GetRatingRankResponse, error) {
	posts, err := h.db.GetScoreRanking(ctx, req.Page, req.PageSize)
	if err != nil {
		return &post.GetRatingRankResponse{
			Code:    constants.PostRatingRankFailCode,
			Message: fmt.Sprintf("failed to get rating rank: %s", err),
		}, nil
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
		Code:  constants.SuccessCode,
		Posts: postList,
		Total: int32(len(postList)),
	}, nil
}

// UpdatePost 更新帖子
func (h *PostHandler) UpdatePost(ctx context.Context, req *post.UpdatePostRequest) (*post.UpdatePostResponse, error) {
	// 参数验证
	if req.PostId == "" {
		return &post.UpdatePostResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to update post: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	// TODO: 实现更新帖子逻辑
	return &post.UpdatePostResponse{
		Code:    constants.SuccessCode,
		Message: "更新成功",
	}, nil
}

// DeletePost 删除帖子
func (h *PostHandler) DeletePost(ctx context.Context, req *post.DeletePostRequest) (*post.DeletePostResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Infof("DeletePost req: %v", req)

	// 参数验证
	if req.PostId == "" {
		return &post.DeletePostResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to delete post: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	err := h.db.DeletePost(ctx, req.PostId)
	if err != nil {
		logger.Errorf("DeletePost failed: %s", err)
		if strings.Contains(fmt.Sprintf("%s", err), "not found") {
			return &post.DeletePostResponse{
				Code:    constants.PostNotFoundCode,
				Message: fmt.Sprintf("failed to delete post: %s", constants.GetErrorMessage(constants.PostNotFoundCode)),
			}, nil
		}
		return &post.DeletePostResponse{
			Code:    constants.PostDeleteFailCode,
			Message: fmt.Sprintf("failed to delete post: %s", err),
		}, nil
	}
	return &post.DeletePostResponse{
		Code: constants.SuccessCode,
	}, nil
}

// GetRecommendPosts 获取推荐帖子
func (h *PostHandler) GetRecommendPosts(ctx context.Context, req *post.GetRecommendPostsRequest) (*post.GetRecommendPostsResponse, error) {
	// TODO: 实现获取推荐帖子逻辑
	return &post.GetRecommendPostsResponse{
		Code:    constants.SuccessCode,
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
		Code:    constants.SuccessCode,
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
		Code:    constants.SuccessCode,
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
		Code:    constants.SuccessCode,
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
		Code:    constants.SuccessCode,
		Message: "获取成功",
		Posts:   []*post.Post{},
		Total:   0,
		HasMore: false,
	}, nil
}

// SearchPosts 搜索帖子
func (h *PostHandler) SearchPosts(ctx context.Context, req *post.SearchPostsRequest) (*post.SearchPostsResponse, error) {
	// 参数验证
	if req.Keyword == "" {
		return &post.SearchPostsResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to search posts: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	// TODO: 实现搜索帖子逻辑
	return &post.SearchPostsResponse{
		Code:    constants.SuccessCode,
		Message: "搜索成功",
		Posts:   []*post.Post{},
		Total:   0,
	}, nil
}

// GetTopic 获取话题详情
func (h *PostHandler) GetTopic(ctx context.Context, req *post.GetTopicRequest) (*post.GetTopicResponse, error) {
	// 参数验证
	if req.TopicId == "" {
		return &post.GetTopicResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to get topic: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	// TODO: 实现获取话题详情逻辑
	return &post.GetTopicResponse{
		Code:    constants.SuccessCode,
		Message: "获取成功",
	}, nil
}

// GetHotTopics 获取热门话题
func (h *PostHandler) GetHotTopics(ctx context.Context, req *post.GetHotTopicsRequest) (*post.GetHotTopicsResponse, error) {
	// TODO: 实现获取热门话题逻辑
	return &post.GetHotTopicsResponse{
		Code:    constants.SuccessCode,
		Message: "获取成功",
		Topics:  []*post.Topic{},
	}, nil
}

// GetTopicCategories 获取话题分类
func (h *PostHandler) GetTopicCategories(ctx context.Context, req *post.GetTopicCategoriesRequest) (*post.GetTopicCategoriesResponse, error) {
	// TODO: 实现获取话题分类逻辑
	return &post.GetTopicCategoriesResponse{
		Code:    constants.SuccessCode,
		Message: "获取成功",
		Topics:  []*post.Topic{},
		Total:   0,
		HasMore: false,
	}, nil
}

// SearchTopics 搜索话题
func (h *PostHandler) SearchTopics(ctx context.Context, req *post.SearchTopicsRequest) (*post.SearchTopicsResponse, error) {
	// 参数验证
	if req.Keyword == "" {
		return &post.SearchTopicsResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to search topics: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	// TODO: 实现搜索话题逻辑
	return &post.SearchTopicsResponse{
		Code:    constants.SuccessCode,
		Message: "搜索成功",
		Topics:  []*post.Topic{},
		Total:   0,
		HasMore: false,
	}, nil
}

// ShareTopic 分享话题
func (h *PostHandler) ShareTopic(ctx context.Context, req *post.ShareTopicRequest) (*post.ShareTopicResponse, error) {
	// 参数验证
	if req.TopicId == "" {
		return &post.ShareTopicResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to share topic: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	// TODO: 实现分享话题逻辑
	return &post.ShareTopicResponse{
		Code:    constants.SuccessCode,
		Message: "分享成功",
	}, nil
}

// GetUserRating 获取用户评分
func (h *PostHandler) GetUserRating(ctx context.Context, req *post.GetUserRatingRequest) (*post.GetUserRatingResponse, error) {
	// 参数验证
	if req.UserId == "" || req.PostId == "" {
		return &post.GetUserRatingResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to get user rating: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	// 获取用户评分
	rating, err := h.db.GetUserRating(ctx, req.UserId, req.PostId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &post.GetUserRatingResponse{
				Code:    constants.PostRatingNotFoundCode,
				Message: "用户未对该帖子评分",
				IsRated: false,
			}, nil
		}
		return &post.GetUserRatingResponse{
			Code:    constants.PostRatingFailCode,
			Message: fmt.Sprintf("failed to get user rating: %s", err),
			IsRated: false,
		}, nil
	}

	return &post.GetUserRatingResponse{
		Code:    constants.SuccessCode,
		Message: "获取成功",
		IsRated: true,
		Score:   utils.Float64Ptr(float64(rating.Score)),
		Comment: rating.Comment,
	}, nil
}

// UpdateRating 更新评分
func (h *PostHandler) UpdateRating(ctx context.Context, req *post.UpdateRatingRequest) (*post.UpdateRatingResponse, error) {
	// 参数验证
	if req.UserId == "" || req.PostId == "" || req.Score < 1 || req.Score > 5 {
		return &post.UpdateRatingResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to update rating: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	// 检查帖子是否存在并获取帖子信息
	postInfo, err := h.db.GetPost(ctx, req.PostId)
	if err != nil {
		return &post.UpdateRatingResponse{
			Code:    constants.PostNotFoundCode,
			Message: fmt.Sprintf("failed to update rating: %s", constants.GetErrorMessage(constants.PostNotFoundCode)),
		}, nil
	}

	// 防止用户给自己的帖子评分
	if postInfo.UserID == req.UserId {
		return &post.UpdateRatingResponse{
			Code:    constants.ValidationErrorCode,
			Message: "不能给自己的帖子评分",
		}, nil
	}

	// 检查用户是否已经评分过
	_, err = h.db.GetUserRating(ctx, req.UserId, req.PostId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &post.UpdateRatingResponse{
				Code:    constants.PostRatingNotFoundCode,
				Message: "用户未对该帖子评分，无法更新",
			}, nil
		}
		return &post.UpdateRatingResponse{
			Code:    constants.PostRatingFailCode,
			Message: fmt.Sprintf("failed to check existing rating: %s", err),
		}, nil
	}

	// 更新评分
	err = h.db.UpdateRating(ctx, req.UserId, req.PostId, int32(req.Score), req.Comment)
	if err != nil {
		return &post.UpdateRatingResponse{
			Code:    constants.PostRatingFailCode,
			Message: fmt.Sprintf("failed to update rating: %s", err),
		}, nil
	}

	// 获取更新后的评分统计信息
	averageScore, totalRatings, err := h.db.GetPostRatingStats(ctx, req.PostId)
	if err != nil {
		return &post.UpdateRatingResponse{
			Code:    constants.PostRatingFailCode,
			Message: fmt.Sprintf("failed to get rating stats: %s", err),
		}, nil
	}

	return &post.UpdateRatingResponse{
		Code:         constants.SuccessCode,
		Message:      "更新成功",
		AverageScore: averageScore,
		TotalRatings: totalRatings,
	}, nil
}

// DeleteRating 删除评分
func (h *PostHandler) DeleteRating(ctx context.Context, req *post.DeleteRatingRequest) (*post.DeleteRatingResponse, error) {
	// 参数验证
	if req.UserId == "" || req.PostId == "" {
		return &post.DeleteRatingResponse{
			Code:    constants.ValidationErrorCode,
			Message: fmt.Sprintf("failed to delete rating: %s", constants.GetErrorMessage(constants.ValidationErrorCode)),
		}, nil
	}

	// 检查用户是否已经评分过
	_, err := h.db.GetUserRating(ctx, req.UserId, req.PostId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &post.DeleteRatingResponse{
				Code:    constants.PostRatingNotFoundCode,
				Message: "用户未对该帖子评分，无法删除",
			}, nil
		}
		return &post.DeleteRatingResponse{
			Code:    constants.PostRatingFailCode,
			Message: fmt.Sprintf("failed to check existing rating: %s", err),
		}, nil
	}

	// 删除评分
	err = h.db.DeleteRating(ctx, req.UserId, req.PostId)
	if err != nil {
		return &post.DeleteRatingResponse{
			Code:    constants.PostRatingFailCode,
			Message: fmt.Sprintf("failed to delete rating: %s", err),
		}, nil
	}

	// 获取删除后的评分统计信息
	averageScore, totalRatings, err := h.db.GetPostRatingStats(ctx, req.PostId)
	if err != nil {
		return &post.DeleteRatingResponse{
			Code:    constants.PostRatingFailCode,
			Message: fmt.Sprintf("failed to get rating stats: %s", err),
		}, nil
	}

	return &post.DeleteRatingResponse{
		Code:         constants.SuccessCode,
		Message:      "删除成功",
		AverageScore: averageScore,
		TotalRatings: totalRatings,
	}, nil
}
