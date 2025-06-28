package post

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/kitex_gen/post"
)

// 评分帖子
func RatePost(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := GetUserID(c)
	if !exists {
		RespondUnauthorized(c)
		return
	}

	// 解析请求参数
	var req post.RatePostRequest
	if !BindAndValidateRequest(c, &req) {
		return
	}

	// 设置用户ID
	req.UserId = userID

	// 调用帖子服务
	resp, err := postClient.RatePost(ctx, &req)
	if err != nil {
		RespondInternalError(c, MsgRatePostFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取用户对帖子的评分
func GetUserRating(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := GetUserID(c)
	if !exists {
		RespondUnauthorized(c)
		return
	}

	// 获取帖子ID参数
	postID, valid := ValidateRequiredPathParam(c, "postId", MsgPostIDEmpty)
	if !valid {
		return
	}

	// 构建请求
	req := &post.GetUserRatingRequest{
		UserId: userID,
		PostId: postID,
	}

	// 调用帖子服务
	resp, err := postClient.GetUserRating(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgGetRatingFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 更新评分
func UpdateRating(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := GetUserID(c)
	if !exists {
		RespondUnauthorized(c)
		return
	}

	// 获取帖子ID参数
	postID, valid := ValidateRequiredPathParam(c, "postId", MsgPostIDEmpty)
	if !valid {
		return
	}

	// 解析请求参数
	var req post.UpdateRatingRequest
	if !BindAndValidateRequest(c, &req) {
		return
	}

	// 设置用户ID和帖子ID
	req.UserId = userID
	req.PostId = postID

	// 调用帖子服务
	resp, err := postClient.UpdateRating(ctx, &req)
	if err != nil {
		RespondInternalError(c, MsgUpdateRatingFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 删除评分
func DeleteRating(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 获取用户ID
	userID, exists := GetUserID(c)
	if !exists {
		RespondUnauthorized(c)
		return
	}

	// 获取帖子ID参数
	postID, valid := ValidateRequiredPathParam(c, "postId", MsgPostIDEmpty)
	if !valid {
		return
	}

	// 构建请求
	req := &post.DeleteRatingRequest{
		UserId: userID,
		PostId: postID,
	}

	// 调用帖子服务
	resp, err := postClient.DeleteRating(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgDeleteRatingFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取评分排行榜
func GetRatingRank(ctx context.Context, c *app.RequestContext) {
	postClient := handler.GetPostClient()
	// 解析分页参数
	page, pageSize := ParsePaginationParamsInt32(c)

	// 解析排行榜类型和日期参数
	rankType := c.Query(ParamRankType)
	if rankType == "" {
		rankType = RankTypeDailyHigh // 默认为每日高分榜
	}
	date := ParseOptionalStringParam(c, ParamDate)

	// 构建请求
	req := &post.GetRatingRankRequest{
		Page:     page,
		PageSize: pageSize,
		RankType: rankType,
		Date:     date,
	}

	// 调用帖子服务
	resp, err := postClient.GetRatingRank(ctx, req)
	if err != nil {
		RespondInternalError(c, MsgGetRankFailed, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
