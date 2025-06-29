package post

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/post"
	"hupu/shared/constants"
)

// 创建帖子
func CreatePost(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 解析请求参数
	var req post.CreatePostRequest
	if !common.BindAndValidateRequest(c, &req) {
		return
	}

	// 设置用户ID
	req.UserId = userID

	// 调用帖子服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.CreatePost(ctx, &req)
	}), "CreatePost", constants.MsgCreatePostFailed)
}

// 获取帖子详情
func GetPost(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取帖子ID参数
	postID, valid := common.ValidateRequiredPathParam(c, "id", constants.MsgPostIDEmpty)
	if !valid {
		return
	}

	// 调用帖子服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.GetPost(ctx, &post.GetPostRequest{PostId: postID})
	}), "GetPost", constants.MsgGetPostFailed)
}

// 更新帖子
func UpdatePost(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 获取帖子ID参数
	postID, valid := common.ValidateRequiredPathParam(c, "id", "帖子ID")
	if !valid {
		return
	}

	// 解析请求参数
	var req post.UpdatePostRequest
	if !common.BindAndValidateRequest(c, &req) {
		return
	}

	// 设置用户ID和帖子ID
	req.UserId = userID
	req.PostId = postID

	// 调用帖子服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.UpdatePost(ctx, &req)
	}), "UpdatePost", constants.MsgUpdatePostFailed)
}

// 删除帖子
func DeletePost(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 获取帖子ID参数
	postID, valid := common.ValidateRequiredPathParam(c, "id", constants.MsgPostIDEmpty)
	if !valid {
		return
	}

	// 构建请求
	req := &post.DeletePostRequest{
		UserId: userID,
		PostId: postID,
	}

	// 调用帖子服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.DeletePost(ctx, req)
	}), "DeletePost", constants.MsgDeletePostFailed)
}

// 获取帖子列表
func GetPostList(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 解析其他查询参数
	// category := common.ParseOptionalStringParam(c, constants.ParamCategory)

	// 构建请求
	req := &post.GetPostListRequest{
		Page:     int64(page),
		PageSize: int64(pageSize),
	}

	// 调用帖子服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.GetPostList(ctx, req)
	}), "GetPostList", constants.MsgGetPostListFailed)
}

// 获取推荐帖子
func GetRecommendPosts(ctx context.Context, c *app.RequestContext) {
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
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.GetRecommendPosts(ctx, req)
	}), "GetRecommendPosts", constants.MsgGetPostListFailed)
}

// 获取热门帖子
func GetHotPosts(ctx context.Context, c *app.RequestContext) {
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
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.GetHotPosts(ctx, req)
	}), "GetHotPosts", constants.MsgGetPostListFailed)
}

// 获取高分帖子
func GetHighScorePosts(ctx context.Context, c *app.RequestContext) {
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
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.GetHighScorePosts(ctx, req)
	}), "GetHighScorePosts", constants.MsgGetPostListFailed)
}

// 获取低分帖子
func GetLowScorePosts(ctx context.Context, c *app.RequestContext) {
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
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.GetLowScorePosts(ctx, req)
	}), "GetLowScorePosts", constants.MsgGetPostListFailed)
}

// 获取争议帖子
func GetControversialPosts(ctx context.Context, c *app.RequestContext) {
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
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.GetControversialPosts(ctx, req)
	}), "GetControversialPosts", constants.MsgGetPostListFailed)
}

// 搜索帖子
func SearchPosts(ctx context.Context, c *app.RequestContext) {
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
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.SearchPosts(ctx, req)
	}), "SearchPosts", constants.MsgSearchPostsFailed)
}
