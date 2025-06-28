package post

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/kitex_gen/post"
)

// 创建帖子
func CreatePost(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := GetUserID(c)
	if !exists {
		RespondUnauthorized(c)
		return
	}

	// 解析请求参数
	var req post.CreatePostRequest
	if !BindAndValidateRequest(c, &req) {
		return
	}

	// 设置用户ID
	req.UserId = userID

	// 调用帖子服务
	resp, err := postClient.CreatePost(ctx, &req)
	if err != nil {
		RespondInternalError(c, MsgCreatePostFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取帖子详情
func GetPost(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取帖子ID参数
	postID, valid := ValidateRequiredPathParam(c, "id", MsgPostIDEmpty)
	if !valid {
		return
	}

	// 调用帖子服务
	req := &post.GetPostRequest{PostId: postID}
	resp, err := postClient.GetPost(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgGetPostFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 更新帖子
func UpdatePost(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := GetUserID(c)
	if !exists {
		RespondUnauthorized(c)
		return
	}

	// 获取帖子ID参数
	postID, valid := ValidateRequiredPathParam(c, "id", MsgPostIDEmpty)
	if !valid {
		return
	}

	// 解析请求参数
	var req post.UpdatePostRequest
	if !BindAndValidateRequest(c, &req) {
		return
	}

	// 设置用户ID和帖子ID
	req.UserId = userID
	req.PostId = postID

	// 调用帖子服务
	resp, err := postClient.UpdatePost(ctx, &req)
	if err != nil {
		RespondInternalError(c, MsgUpdatePostFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 删除帖子
func DeletePost(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := GetUserID(c)
	if !exists {
		RespondUnauthorized(c)
		return
	}

	// 获取帖子ID参数
	postID, valid := ValidateRequiredPathParam(c, "id", MsgPostIDEmpty)
	if !valid {
		return
	}

	// 构建请求
	req := &post.DeletePostRequest{
		UserId: userID,
		PostId: postID,
	}

	// 调用帖子服务
	resp, err := postClient.DeletePost(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgDeletePostFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取帖子列表
func GetPostList(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := ParsePaginationParams(c)

	// 解析其他查询参数
	userIDStr := ParseOptionalStringParam(c, ParamUserID)
	topicIDStr := ParseOptionalStringParam(c, ParamTopicID)
	categoryStr := c.Query(ParamCategory)
	sortType := ParseOptionalStringParam(c, ParamSortType)
	isAnonymous := ParseOptionalBoolParam(c, ParamIsAnonymous)

	// 构建请求
	req := &post.GetPostListRequest{
		Page:     page,
		PageSize: pageSize,
	}

	if userIDStr != nil {
		req.UserId = userIDStr
	}

	if topicIDStr != nil {
		req.TopicId = topicIDStr
	}

	if categoryStr != "" {
		if category, err := strconv.ParseInt(categoryStr, 10, 32); err == nil {
			categoryEnum := post.PostCategory(category)
			req.Category = &categoryEnum
		}
	}

	if sortType != nil {
		req.SortType = sortType
	}

	if isAnonymous != nil {
		req.IsAnonymous = isAnonymous
	}

	// 调用帖子服务
	resp, err := postClient.GetPostList(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgGetPostListFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取推荐帖子
func GetRecommendPosts(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := ParsePaginationParamsInt32(c)

	// 解析其他查询参数
	category := ParseOptionalStringParam(c, ParamCategory)
	tag := ParseOptionalStringParam(c, ParamTag)

	// 构建请求
	req := &post.GetRecommendPostsRequest{
		Page:     page,
		PageSize: pageSize,
		Category: category,
		Tag:      tag,
	}

	// 调用帖子服务
	resp, err := postClient.GetRecommendPosts(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgGetPostListFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取热门帖子
func GetHotPosts(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := ParsePaginationParamsInt32(c)

	// 解析其他查询参数
	category := ParseOptionalStringParam(c, ParamCategory)
	tag := ParseOptionalStringParam(c, ParamTag)

	// 构建请求
	req := &post.GetHotPostsRequest{
		Page:     page,
		PageSize: pageSize,
		Category: category,
		Tag:      tag,
	}

	// 调用帖子服务
	resp, err := postClient.GetHotPosts(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgGetPostListFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取高分帖子
func GetHighScorePosts(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := ParsePaginationParamsInt32(c)

	// 解析其他查询参数
	category := ParseOptionalStringParam(c, ParamCategory)
	tag := ParseOptionalStringParam(c, ParamTag)

	// 构建请求
	req := &post.GetHighScorePostsRequest{
		Page:     page,
		PageSize: pageSize,
		Category: category,
		Tag:      tag,
	}

	// 调用帖子服务
	resp, err := postClient.GetHighScorePosts(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgGetPostListFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取低分帖子
func GetLowScorePosts(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := ParsePaginationParamsInt32(c)

	// 解析其他查询参数
	category := ParseOptionalStringParam(c, ParamCategory)
	tag := ParseOptionalStringParam(c, ParamTag)

	// 构建请求
	req := &post.GetLowScorePostsRequest{
		Page:     page,
		PageSize: pageSize,
		Category: category,
		Tag:      tag,
	}

	// 调用帖子服务
	resp, err := postClient.GetLowScorePosts(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgGetPostListFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取争议帖子
func GetControversialPosts(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := ParsePaginationParamsInt32(c)

	// 解析其他查询参数
	category := ParseOptionalStringParam(c, ParamCategory)
	tag := ParseOptionalStringParam(c, ParamTag)

	// 构建请求
	req := &post.GetControversialPostsRequest{
		Page:     page,
		PageSize: pageSize,
		Category: category,
		Tag:      tag,
	}

	// 调用帖子服务
	resp, err := postClient.GetControversialPosts(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgGetPostListFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 搜索帖子
func SearchPosts(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取搜索关键词
	keyword := c.Query(ParamKeyword)
	if keyword == "" {
		RespondBadRequest(c, MsgKeywordEmpty)
		return
	}

	// 解析分页参数
	page, pageSize := ParsePaginationParamsInt32(c)

	// 构建请求
	req := &post.SearchPostsRequest{
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	// 调用帖子服务
	resp, err := postClient.SearchPosts(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgSearchPostsFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
