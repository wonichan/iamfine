package post

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/post"
	"hupu/shared/constants"
)

// 评分帖子
func RatePost(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 解析请求参数
	var req post.RatePostRequest
	if !common.BindAndValidateRequest(c, &req) {
		return
	}

	// 设置用户ID
	req.UserId = userID

	// 调用帖子服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.RatePost(ctx, &req)
	}), "RatePost", constants.MsgRatePostFailed)
}

// 获取用户对帖子的评分
func GetUserRating(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 获取帖子ID参数
	postID, valid := common.ValidateRequiredPathParam(c, "postId", constants.MsgPostIDEmpty)
	if !valid {
		return
	}

	// 构建请求
	req := &post.GetUserRatingRequest{
		UserId: userID,
		PostId: postID,
	}

	// 调用帖子服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.GetUserRating(ctx, req)
	}), "GetUserRating", constants.MsgGetRatingFailed)
}

// 更新评分
func UpdateRating(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 获取帖子ID参数
	postID, valid := common.ValidateRequiredPathParam(c, "postId", constants.MsgPostIDEmpty)
	if !valid {
		return
	}

	// 解析请求参数
	var req post.UpdateRatingRequest
	if !common.BindAndValidateRequest(c, &req) {
		return
	}

	// 设置用户ID和帖子ID
	req.UserId = userID
	req.PostId = postID

	// 调用帖子服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.UpdateRating(ctx, &req)
	}), "UpdateRating", constants.MsgUpdateRatingFailed)
}

// 删除评分
func DeleteRating(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 获取帖子ID参数
	postID, valid := common.ValidateRequiredPathParam(c, "postId", constants.MsgPostIDEmpty)
	if !valid {
		return
	}

	// 构建请求
	req := &post.DeleteRatingRequest{
		UserId: userID,
		PostId: postID,
	}

	// 调用帖子服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.DeleteRating(ctx, req)
	}), "DeleteRating", constants.MsgDeleteRatingFailed)
}

// 获取评分排行榜
func GetRatingRank(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 解析排行榜类型和日期参数
	rankType := c.Query(constants.ParamRankType)
	if rankType == "" {
		rankType = constants.RankTypeDailyHigh // 默认为每日高分榜
	}
	date := common.ParseOptionalStringParam(c, constants.ParamDate)

	// 构建请求
	req := &post.GetRatingRankRequest{
		Page:     page,
		PageSize: pageSize,
		RankType: rankType,
		Date:     date,
	}

	// 调用帖子服务
	common.CallService(c, common.ServiceCall(func() (any, error) {
		return postClient.GetRatingRank(ctx, req)
	}), "GetRatingRank", constants.MsgGetRankFailed)
}
