package post

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/post"
	"hupu/shared/constants"
	"hupu/shared/log"
)

// 创建帖子
func CreatePost(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		logger.Errorf("CreatePost GetUserIDFromContext failed")
		return
	}

	// 解析请求参数
	var req post.CreatePostRequest
	if !common.BindAndValidateRequest(c, &req) {
		logger.Errorf("CreatePost BindAndValidateRequest failed")
		return
	}

	// 设置用户ID
	req.UserId = userID

	// 调用帖子服务
	resp, err := postClient.CreatePost(ctx, &req)
	if err != nil {
		common.HandleRpcError(c, "CreatePost", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "CreatePost", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// 获取帖子详情
func GetPost(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	postClient := handler.GetPostClient()
	// 获取帖子ID参数
	postID, valid := common.ValidateRequiredPathParam(c, "id", constants.MsgPostIDEmpty)
	if !valid {
		logger.Errorf("GetPost ValidateRequiredPathParam failed")
		return
	}

	// 调用帖子服务
	resp, err := postClient.GetPost(ctx, &post.GetPostRequest{PostId: postID})
	if err != nil {
		common.HandleRpcError(c, "GetPost", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetPost", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// 更新帖子
func UpdatePost(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		logger.Errorf("UpdatePost GetUserIDFromContext failed")
		return
	}

	// 获取帖子ID参数
	postID, valid := common.ValidateRequiredPathParam(c, "id", "帖子ID")
	if !valid {
		logger.Errorf("UpdatePost ValidateRequiredPathParam failed")
		return
	}

	// 解析请求参数
	var req post.UpdatePostRequest
	if !common.BindAndValidateRequest(c, &req) {
		logger.Errorf("UpdatePost BindAndValidateRequest failed")
		return
	}

	// 设置用户ID和帖子ID
	req.UserId = userID
	req.PostId = postID

	// 调用帖子服务
	resp, err := postClient.UpdatePost(ctx, &req)
	if err != nil {
		common.HandleRpcError(c, "UpdatePost", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "UpdatePost", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// 删除帖子
func DeletePost(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		logger.Errorf("DeletePost GetUserIDFromContext failed")
		return
	}

	// 获取帖子ID参数
	postID, valid := common.ValidateRequiredPathParam(c, "id", constants.MsgPostIDEmpty)
	if !valid {
		logger.Errorf("DeletePost ValidateRequiredPathParam failed")
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
		common.HandleRpcError(c, "DeletePost", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "DeletePost", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// 获取帖子列表
func GetPostList(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 解析其他查询参数
	category := int64(common.ParseOptionalIntParam(c, constants.ParamCategory))
	topic := common.ParseOptionalStringParam(c, constants.ParamTopicID)
	sortType := common.ParseOptionalStringParam(c, constants.ParamSortType)

	// 构建请求
	req := &post.GetPostListRequest{
		Page:     int64(page),
		PageSize: int64(pageSize),
		Category: (*post.PostCategory)(&category),
		TopicId:  topic,
		SortType: sortType,
	}

	// 调用帖子服务
	resp, err := handler.GetPostClient().GetPostList(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "GetPostList", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetPostList", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// 获取推荐帖子
func GetRecommendPosts(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 解析其他查询参数
	category := common.ParseOptionalStringParam(c, constants.ParamCategory)
	tag := common.ParseOptionalStringParam(c, constants.ParamTag)

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
		common.HandleRpcError(c, "GetRecommendPosts", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetRecommendPosts", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// 获取热门帖子
func GetHotPosts(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 解析其他查询参数
	category := common.ParseOptionalStringParam(c, constants.ParamCategory)
	tag := common.ParseOptionalStringParam(c, constants.ParamTag)

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
		common.HandleRpcError(c, "GetHotPosts", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetHotPosts", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// 获取高分帖子
func GetHighScorePosts(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 解析其他查询参数
	category := common.ParseOptionalStringParam(c, constants.ParamCategory)
	tag := common.ParseOptionalStringParam(c, constants.ParamTag)

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
		common.HandleRpcError(c, "GetHighScorePosts", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetHighScorePosts", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// 获取低分帖子
func GetLowScorePosts(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 解析其他查询参数
	category := common.ParseOptionalStringParam(c, constants.ParamCategory)
	tag := common.ParseOptionalStringParam(c, constants.ParamTag)

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
		common.HandleRpcError(c, "GetLowScorePosts", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetLowScorePosts", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// 获取争议帖子
func GetControversialPosts(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 解析其他查询参数
	category := common.ParseOptionalStringParam(c, constants.ParamCategory)
	tag := common.ParseOptionalStringParam(c, constants.ParamTag)

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
		common.HandleRpcError(c, "GetControversialPosts", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetControversialPosts", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// 搜索帖子
func SearchPosts(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	postClient := handler.GetPostClient()
	// 获取搜索关键词
	keyword := c.Query(constants.ParamKeyword)
	if keyword == "" {
		common.RespondBadRequest(c, constants.MsgKeywordEmpty)
		return
	}

	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 构建请求
	req := &post.SearchPostsRequest{
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	// 调用帖子服务
	resp, err := postClient.SearchPosts(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "SearchPosts", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "SearchPosts", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}
