package post

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/kitex_gen/post"
)

// 创建话题
func CreateTopic(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 解析请求参数
	var req post.CreateTopicRequest
	if !BindAndValidateRequest(c, &req) {
		return
	}

	// 调用帖子服务
	resp, err := postClient.CreateTopic(ctx, &req)
	if err != nil {
		RespondInternalError(c, MsgCreateTopicFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取话题详情
func GetTopic(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取话题ID参数
	topicID, valid := ValidateRequiredPathParam(c, "id", MsgTopicIDEmpty)
	if !valid {
		return
	}

	// 调用帖子服务
	req := &post.GetTopicRequest{TopicId: topicID}
	resp, err := postClient.GetTopic(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgGetTopicFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取话题列表
func GetTopicList(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := ParsePaginationParamsInt32(c)

	// 解析排序类型
	sortType := ParseOptionalStringParam(c, ParamSortType)

	// 构建请求
	req := &post.GetTopicListRequest{
		Page:     page,
		PageSize: pageSize,
		SortType: sortType,
	}

	// 调用帖子服务
	resp, err := postClient.GetTopicList(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgGetTopicListFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取热门话题
func GetHotTopics(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 解析限制数量参数
	limit := ParseOptionalIntParam(c, ParamLimit)

	// 构建请求
	req := &post.GetHotTopicsRequest{
		Limit: limit,
	}

	// 调用帖子服务
	resp, err := postClient.GetHotTopics(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgGetTopicListFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取话题分类
func GetTopicCategories(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 解析限制数量参数
	limit := ParseOptionalIntParam(c, ParamLimit)

	// 构建请求
	req := &post.GetTopicCategoriesRequest{
		Limit: limit,
	}

	// 调用帖子服务
	resp, err := postClient.GetTopicCategories(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgGetTopicListFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 搜索话题
func SearchTopics(ctx context.Context, c *app.RequestContext) {
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
	req := &post.SearchTopicsRequest{
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	// 调用帖子服务
	resp, err := postClient.SearchTopics(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgSearchTopicsFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 分享话题
func ShareTopic(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := GetUserID(c)
	if !exists {
		RespondUnauthorized(c)
		return
	}

	// 获取话题ID参数
	topicID, valid := ValidateRequiredPathParam(c, "id", MsgTopicIDEmpty)
	if !valid {
		return
	}

	// 构建请求
	req := &post.ShareTopicRequest{
		TopicId: topicID,
		UserId:  userID,
	}

	// 调用帖子服务
	resp, err := postClient.ShareTopic(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgShareTopicFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
