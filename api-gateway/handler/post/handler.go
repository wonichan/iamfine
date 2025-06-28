package post

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// ==================== 帖子管理Handler ====================

// CreatePostHandler 创建帖子
func CreatePostHandler(ctx context.Context, c *app.RequestContext) {
	CreatePost(ctx, c)
}

// GetPostHandler 获取帖子详情
func GetPostHandler(ctx context.Context, c *app.RequestContext) {
	GetPost(ctx, c)
}

// UpdatePostHandler 更新帖子
func UpdatePostHandler(ctx context.Context, c *app.RequestContext) {
	UpdatePost(ctx, c)
}

// DeletePostHandler 删除帖子
func DeletePostHandler(ctx context.Context, c *app.RequestContext) {
	DeletePost(ctx, c)
}

// GetPostListHandler 获取帖子列表
func GetPostListHandler(ctx context.Context, c *app.RequestContext) {
	GetPostList(ctx, c)
}

// GetRecommendPostsHandler 获取推荐帖子
func GetRecommendPostsHandler(ctx context.Context, c *app.RequestContext) {
	GetRecommendPosts(ctx, c)
}

// GetHotPostsHandler 获取热门帖子
func GetHotPostsHandler(ctx context.Context, c *app.RequestContext) {
	GetHotPosts(ctx, c)
}

// GetHighScorePostsHandler 获取高分帖子
func GetHighScorePostsHandler(ctx context.Context, c *app.RequestContext) {
	GetHighScorePosts(ctx, c)
}

// GetLowScorePostsHandler 获取低分帖子
func GetLowScorePostsHandler(ctx context.Context, c *app.RequestContext) {
	GetLowScorePosts(ctx, c)
}

// GetControversialPostsHandler 获取争议帖子
func GetControversialPostsHandler(ctx context.Context, c *app.RequestContext) {
	GetControversialPosts(ctx, c)
}

// SearchPostsHandler 搜索帖子
func SearchPostsHandler(ctx context.Context, c *app.RequestContext) {
	SearchPosts(ctx, c)
}

// ==================== 话题管理Handler ====================

// CreateTopicHandler 创建话题
func CreateTopicHandler(ctx context.Context, c *app.RequestContext) {
	CreateTopic(ctx, c)
}

// GetTopicHandler 获取话题详情
func GetTopicHandler(ctx context.Context, c *app.RequestContext) {
	GetTopic(ctx, c)
}

// GetTopicListHandler 获取话题列表
func GetTopicListHandler(ctx context.Context, c *app.RequestContext) {
	GetTopicList(ctx, c)
}

// GetHotTopicsHandler 获取热门话题
func GetHotTopicsHandler(ctx context.Context, c *app.RequestContext) {
	GetHotTopics(ctx, c)
}

// GetTopicCategoriesHandler 获取话题分类
func GetTopicCategoriesHandler(ctx context.Context, c *app.RequestContext) {
	GetTopicCategories(ctx, c)
}

// SearchTopicsHandler 搜索话题
func SearchTopicsHandler(ctx context.Context, c *app.RequestContext) {
	SearchTopics(ctx, c)
}

// ShareTopicHandler 分享话题
func ShareTopicHandler(ctx context.Context, c *app.RequestContext) {
	ShareTopic(ctx, c)
}

// ==================== 收藏管理Handler ====================

// CollectPostHandler 收藏帖子
func CollectPostHandler(ctx context.Context, c *app.RequestContext) {
	CollectPost(ctx, c)
}

// UncollectPostHandler 取消收藏帖子
func UncollectPostHandler(ctx context.Context, c *app.RequestContext) {
	UncollectPost(ctx, c)
}

// GetCollectedPostsHandler 获取收藏的帖子列表
func GetCollectedPostsHandler(ctx context.Context, c *app.RequestContext) {
	GetCollectedPosts(ctx, c)
}

// ==================== 评分管理Handler ====================

// RatePostHandler 评分帖子
func RatePostHandler(ctx context.Context, c *app.RequestContext) {
	RatePost(ctx, c)
}

// GetUserRatingHandler 获取用户对帖子的评分
func GetUserRatingHandler(ctx context.Context, c *app.RequestContext) {
	GetUserRating(ctx, c)
}

// UpdateRatingHandler 更新评分
func UpdateRatingHandler(ctx context.Context, c *app.RequestContext) {
	UpdateRating(ctx, c)
}

// DeleteRatingHandler 删除评分
func DeleteRatingHandler(ctx context.Context, c *app.RequestContext) {
	DeleteRating(ctx, c)
}

// GetRatingRankHandler 获取评分排行榜
func GetRatingRankHandler(ctx context.Context, c *app.RequestContext) {
	GetRatingRank(ctx, c)
}
