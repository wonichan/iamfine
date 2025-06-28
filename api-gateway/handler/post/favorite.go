package post

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/kitex_gen/post"
)

// 收藏帖子
func CollectPost(ctx context.Context, c *app.RequestContext) {
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
	req := &post.CollectPostRequest{
		UserId: userID,
		PostId: postID,
	}

	// 调用帖子服务
	resp, err := postClient.CollectPost(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgCollectPostFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 取消收藏帖子
func UncollectPost(ctx context.Context, c *app.RequestContext) {
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
	req := &post.UncollectPostRequest{
		UserId: userID,
		PostId: postID,
	}

	// 调用帖子服务
	resp, err := postClient.UncollectPost(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgUncollectFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取收藏的帖子列表
func GetCollectedPosts(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := GetUserID(c)
	if !exists {
		RespondUnauthorized(c)
		return
	}

	// 解析分页参数
	page, pageSize := ParsePaginationParamsInt32(c)

	// 构建请求
	req := &post.GetCollectedPostsRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	// 调用帖子服务
	resp, err := postClient.GetCollectedPosts(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgGetCollectedFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
