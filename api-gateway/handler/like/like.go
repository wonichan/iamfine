package like

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/like"
	"hupu/shared/log"
)

// ==================== 帖子点赞实现函数 ====================

// LikePost 点赞帖子
func LikePost(ctx context.Context, c *app.RequestContext) {
	// 获取帖子ID
	postID := c.Param("id")
	if postID == "" {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, MsgTargetIDRequired)
		return
	}

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.ErrorResponseFunc(c, HTTPStatusUnauthorized, CodeError, MsgUnauthorized)
		return
	}

	// 构建请求
	req := &like.LikeRequest{
		UserId:     userID,
		TargetId:   postID,
		TargetType: TargetTypePost,
	}

	// 调用服务
	CallLikeService(ctx, c, func() (interface{}, error) {
		return GetLikeClient().Like(ctx, req)
	}, MsgLikeFailed)
}

// UnlikePost 取消点赞帖子
func UnlikePost(ctx context.Context, c *app.RequestContext) {
	// 获取帖子ID
	postID := c.Param("id")
	if postID == "" {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, MsgTargetIDRequired)
		return
	}

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.ErrorResponseFunc(c, HTTPStatusUnauthorized, CodeError, MsgUnauthorized)
		return
	}

	// 构建请求
	req := &like.UnlikeRequest{
		UserId:     userID,
		TargetId:   postID,
		TargetType: TargetTypePost,
	}

	// 调用服务
	CallLikeService(ctx, c, func() (interface{}, error) {
		return GetLikeClient().Unlike(ctx, req)
	}, MsgUnlikeFailed)
}

// GetPostLikeCount 获取帖子点赞数量
func GetPostLikeCount(ctx context.Context, c *app.RequestContext) {
	// 获取帖子ID
	postID := c.Param("id")
	if postID == "" {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, MsgTargetIDRequired)
		return
	}

	// 构建请求
	req := &like.GetLikeCountRequest{
		TargetId:   postID,
		TargetType: TargetTypePost,
	}

	// 调用服务
	CallLikeService(ctx, c, func() (interface{}, error) {
		return GetLikeClient().GetLikeCount(ctx, req)
	}, MsgGetLikeCountFailed)
}

// CheckPostLikeStatus 检查帖子点赞状态
func CheckPostLikeStatus(ctx context.Context, c *app.RequestContext) {
	// 获取帖子ID
	postID := c.Param("id")
	if postID == "" {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, MsgTargetIDRequired)
		return
	}

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.ErrorResponseFunc(c, HTTPStatusUnauthorized, CodeError, MsgUnauthorized)
		return
	}

	// 构建请求
	req := &like.LikeRequest{
		UserId:     userID,
		TargetId:   postID,
		TargetType: TargetTypePost,
	}

	// 调用服务
	CallLikeService(ctx, c, func() (interface{}, error) {
		return GetLikeClient().IsLiked(ctx, req)
	}, MsgCheckLikeStatusFailed)
}

// ==================== 评论点赞实现函数 ====================

// LikeComment 点赞评论
func LikeComment(ctx context.Context, c *app.RequestContext) {
	// 获取评论ID
	commentID := c.Param("id")
	if commentID == "" {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, MsgTargetIDRequired)
		return
	}

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.ErrorResponseFunc(c, HTTPStatusUnauthorized, CodeError, MsgUnauthorized)
		return
	}

	// 构建请求
	req := &like.LikeRequest{
		UserId:     userID,
		TargetId:   commentID,
		TargetType: TargetTypeComment,
	}

	// 调用服务
	CallLikeService(ctx, c, func() (interface{}, error) {
		return GetLikeClient().Like(ctx, req)
	}, MsgLikeFailed)
}

// UnlikeComment 取消点赞评论
func UnlikeComment(ctx context.Context, c *app.RequestContext) {
	// 获取评论ID
	commentID := c.Param("id")
	if commentID == "" {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, MsgTargetIDRequired)
		return
	}

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.ErrorResponseFunc(c, HTTPStatusUnauthorized, CodeError, MsgUnauthorized)
		return
	}

	// 构建请求
	req := &like.UnlikeRequest{
		UserId:     userID,
		TargetId:   commentID,
		TargetType: TargetTypeComment,
	}

	// 调用服务
	CallLikeService(ctx, c, func() (interface{}, error) {
		return GetLikeClient().Unlike(ctx, req)
	}, MsgUnlikeFailed)
}

// GetCommentLikeCount 获取评论点赞数量
func GetCommentLikeCount(ctx context.Context, c *app.RequestContext) {
	// 获取评论ID
	commentID := c.Param("id")
	if commentID == "" {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, MsgTargetIDRequired)
		return
	}

	// 构建请求
	req := &like.GetLikeCountRequest{
		TargetId:   commentID,
		TargetType: TargetTypeComment,
	}

	// 调用服务
	CallLikeService(ctx, c, func() (interface{}, error) {
		return GetLikeClient().GetLikeCount(ctx, req)
	}, MsgGetLikeCountFailed)
}

// CheckCommentLikeStatus 检查评论点赞状态
func CheckCommentLikeStatus(ctx context.Context, c *app.RequestContext) {
	// 获取评论ID
	commentID := c.Param("id")
	if commentID == "" {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, MsgTargetIDRequired)
		return
	}

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.ErrorResponseFunc(c, HTTPStatusUnauthorized, CodeError, MsgUnauthorized)
		return
	}

	// 构建请求
	req := &like.LikeRequest{
		UserId:     userID,
		TargetId:   commentID,
		TargetType: TargetTypeComment,
	}

	// 调用服务
	CallLikeService(ctx, c, func() (interface{}, error) {
		return GetLikeClient().IsLiked(ctx, req)
	}, MsgCheckLikeStatusFailed)
}

// ==================== 通用点赞实现函数 ====================

// GetLikeList 获取用户点赞列表
func GetLikeList(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.ErrorResponseFunc(c, HTTPStatusUnauthorized, CodeError, MsgUnauthorized)
		return
	}

	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 构建请求
	req := &like.GetLikeListRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	// 调用服务
	CallLikeService(ctx, c, func() (interface{}, error) {
		return GetLikeClient().GetLikeList(ctx, req)
	}, MsgGetLikeListFailed)
}

// IsLiked 检查是否已点赞
func IsLiked(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.ErrorResponseFunc(c, HTTPStatusUnauthorized, CodeError, MsgUnauthorized)
		return
	}

	// 获取查询参数
	targetID := c.Query("target_id")
	targetType := c.Query("target_type")

	// 验证参数
	if err := ValidateTargetParams(targetID, targetType); err != nil {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, err.Error())
		return
	}

	// 构建请求
	req := &like.LikeRequest{
		UserId:     userID,
		TargetId:   targetID,
		TargetType: targetType,
	}

	// 调用服务
	CallLikeService(ctx, c, func() (interface{}, error) {
		return GetLikeClient().IsLiked(ctx, req)
	}, MsgCheckLikeStatusFailed)
}

// GetLikeCount 获取点赞数量
func GetLikeCount(ctx context.Context, c *app.RequestContext) {
	// 获取查询参数
	targetID := c.Query("target_id")
	targetType := c.Query("target_type")

	// 验证参数
	if err := ValidateTargetParams(targetID, targetType); err != nil {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, err.Error())
		return
	}

	// 构建请求
	req := &like.GetLikeCountRequest{
		TargetId:   targetID,
		TargetType: targetType,
	}

	// 调用服务
	CallLikeService(ctx, c, func() (interface{}, error) {
		return GetLikeClient().GetLikeCount(ctx, req)
	}, MsgGetLikeCountFailed)
}

// GetLikeUsers 获取点赞用户列表
func GetLikeUsers(ctx context.Context, c *app.RequestContext) {
	// 获取查询参数
	targetID := c.Query("target_id")
	targetType := c.Query("target_type")

	// 验证参数
	if err := ValidateTargetParams(targetID, targetType); err != nil {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, err.Error())
		return
	}

	// 解析分页参数
	page, pageSize := common.ParsePaginationParams(c)

	// 构建请求
	req := &like.GetLikeUsersRequest{
		TargetId:   targetID,
		TargetType: targetType,
		Page:       page,
		PageSize:   pageSize,
	}

	// 调用服务
	CallLikeService(ctx, c, func() (interface{}, error) {
		return GetLikeClient().GetLikeUsers(ctx, req)
	}, MsgGetLikeUsersFailed)
}

// ==================== 兼容性实现函数 ====================

// Like 兼容旧的点赞接口
func Like(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.ErrorResponseFunc(c, HTTPStatusUnauthorized, CodeError, MsgUnauthorized)
		return
	}

	// 解析请求参数
	var req like.LikeRequest
	if err := c.BindAndValidate(&req); err != nil {
		log.GetLogger().Errorf("Bind request error: %v", err)
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, MsgParamError+": "+err.Error())
		return
	}

	// 设置用户ID
	req.UserId = userID

	// 验证目标类型
	if err := ValidateTargetParams(req.TargetId, req.TargetType); err != nil {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, err.Error())
		return
	}

	// 调用服务
	CallLikeService(ctx, c, func() (interface{}, error) {
		return GetLikeClient().Like(ctx, &req)
	}, MsgLikeFailed)
}

// Unlike 兼容旧的取消点赞接口
func Unlike(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.ErrorResponseFunc(c, HTTPStatusUnauthorized, CodeError, MsgUnauthorized)
		return
	}

	// 解析请求参数
	var req like.UnlikeRequest
	if err := c.BindAndValidate(&req); err != nil {
		log.GetLogger().Errorf("Bind request error: %v", err)
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, MsgParamError+": "+err.Error())
		return
	}

	// 设置用户ID
	req.UserId = userID

	// 验证目标类型
	if err := ValidateTargetParams(req.TargetId, req.TargetType); err != nil {
		common.ErrorResponseFunc(c, HTTPStatusBadRequest, CodeError, err.Error())
		return
	}

	// 调用服务
	CallLikeService(ctx, c, func() (interface{}, error) {
		return GetLikeClient().Unlike(ctx, &req)
	}, MsgUnlikeFailed)
}
