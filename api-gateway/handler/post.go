package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/kitex_gen/post"
	"hupu/shared/log"
)

// 创建帖子
func CreatePost(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 解析请求参数
	var req post.CreatePostRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 设置用户ID
	req.UserId = userID.(string)

	// 调用帖子服务
	resp, err := postClient.CreatePost(ctx, &req)
	if err != nil {
		log.GetLogger().Errorf("CreatePost error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建帖子失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取帖子详情
func GetPost(ctx context.Context, c *app.RequestContext) {
	// 获取帖子ID参数
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "帖子ID不能为空",
		})
		return
	}

	// 调用帖子服务
	req := &post.GetPostRequest{PostId: postID}
	resp, err := postClient.GetPost(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetPost error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取帖子失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取帖子列表
func GetPostList(ctx context.Context, c *app.RequestContext) {
	// 解析查询参数
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	userIDStr := c.Query("user_id")
	topicIDStr := c.Query("topic_id")
	categoryStr := c.Query("category")
	sortType := c.Query("sort_type")
	isAnonymousStr := c.Query("is_anonymous")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 构建请求
	req := &post.GetPostListRequest{
		Page:     page,
		PageSize: pageSize,
	}

	if userIDStr != "" {
		req.UserId = &userIDStr
	}

	if topicIDStr != "" {
		req.TopicId = &topicIDStr
	}

	if categoryStr != "" {
		if category, err := strconv.ParseInt(categoryStr, 10, 32); err == nil {
			categoryEnum := post.PostCategory(category)
			req.Category = &categoryEnum
		}
	}

	if sortType != "" {
		req.SortType = &sortType
	}

	if isAnonymousStr != "" {
		if isAnonymous, err := strconv.ParseBool(isAnonymousStr); err == nil {
			req.IsAnonymous = &isAnonymous
		}
	}

	// 调用帖子服务
	resp, err := postClient.GetPostList(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetPostList error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取帖子列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 创建话题
func CreateTopic(ctx context.Context, c *app.RequestContext) {
	// 解析请求参数
	var req post.CreateTopicRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 调用帖子服务
	resp, err := postClient.CreateTopic(ctx, &req)
	if err != nil {
		log.GetLogger().Errorf("CreateTopic error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建话题失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取话题列表
func GetTopicList(ctx context.Context, c *app.RequestContext) {
	// 解析查询参数
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	sortType := c.Query("sort_type")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 构建请求
	req := &post.GetTopicListRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	if sortType != "" {
		req.SortType = &sortType
	}

	// 调用帖子服务
	resp, err := postClient.GetTopicList(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetTopicList error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取话题列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 收藏帖子
func CollectPost(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取帖子ID参数
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "帖子ID不能为空",
		})
		return
	}

	// 构建请求
	req := &post.CollectPostRequest{
		UserId: userID.(string),
		PostId: postID,
	}

	// 调用帖子服务
	resp, err := postClient.CollectPost(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("CollectPost error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "收藏帖子失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 取消收藏帖子
func UncollectPost(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 获取帖子ID参数
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "帖子ID不能为空",
		})
		return
	}

	// 构建请求
	req := &post.UncollectPostRequest{
		UserId: userID.(string),
		PostId: postID,
	}

	// 调用帖子服务
	resp, err := postClient.UncollectPost(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("UncollectPost error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "取消收藏失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取收藏的帖子列表
func GetCollectedPosts(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 解析查询参数
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 构建请求
	req := &post.GetCollectedPostsRequest{
		UserId:   userID.(string),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	// 调用帖子服务
	resp, err := postClient.GetCollectedPosts(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetCollectedPosts error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取收藏列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 评分帖子
func RatePost(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 解析请求参数
	var req post.RatePostRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 设置用户ID
	req.UserId = userID.(string)

	// 调用帖子服务
	resp, err := postClient.RatePost(ctx, &req)
	if err != nil {
		log.GetLogger().Errorf("RatePost error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "评分失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取评分排行榜
func GetRatingRank(ctx context.Context, c *app.RequestContext) {
	// 解析查询参数
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	rankType := c.Query("rank_type")
	date := c.Query("date")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	if rankType == "" {
		rankType = "daily_high" // 默认为每日高分榜
	}

	// 构建请求
	req := &post.GetRatingRankRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		RankType: rankType,
	}

	if date != "" {
		req.Date = &date
	}

	// 调用帖子服务
	resp, err := postClient.GetRatingRank(ctx, req)
	if err != nil {
		log.GetLogger().Errorf("GetRatingRank error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取排行榜失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
