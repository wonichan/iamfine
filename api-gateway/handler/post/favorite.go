package post

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/post"
	"hupu/shared/constants"
)

// 收藏帖子
func CollectPost(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := common.GetUserID(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 获取帖子ID参数
	postID, valid := common.ValidateRequiredPathParam(c, "id", "帖子ID")
	if !valid {
		return
	}

	// 构建请求
	req := &post.CollectPostRequest{
		UserId: userID,
		PostId: postID,
	}

	// 调用帖子服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.CollectPost(ctx, req)
	}), "CollectPost", constants.MsgCollectPostFailed)
}

// 取消收藏帖子
func UncollectPost(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := common.GetUserID(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 获取帖子ID参数
	postID, valid := common.ValidateRequiredPathParam(c, "id", "帖子ID")
	if !valid {
		return
	}

	// 构建请求
	req := &post.UncollectPostRequest{
		UserId: userID,
		PostId: postID,
	}

	// 调用帖子服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.UncollectPost(ctx, req)
	}), "UncollectPost", constants.MsgUncollectFailed)
}

// 获取收藏的帖子列表
func GetCollectedPosts(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := common.GetUserID(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 构建请求
	req := &post.GetCollectedPostsRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	// 调用帖子服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.GetCollectedPosts(ctx, req)
	}), "GetCollectedPosts", constants.MsgCollectPostFailed)
}
