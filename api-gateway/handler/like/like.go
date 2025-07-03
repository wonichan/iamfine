package like

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/like"
	"hupu/shared/constants"
	"hupu/shared/log"
)

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateTargetParams 验证目标参数
func ValidateTargetParams(targetID, targetType string) error {
	if targetID == "" {
		return ValidationError{Field: "target_id", Message: "目标ID不能为空"}
	}
	if targetType == "" {
		return ValidationError{Field: "target_type", Message: "目标类型不能为空"}
	}
	if targetType != constants.TargetTypePost && targetType != constants.TargetTypeComment {
		return ValidationError{Field: "target_type", Message: "无效的目标类型"}
	}
	return nil
}

// ==================== 帖子点赞实现函数 ====================

// LikePost 点赞帖子
func LikePost(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	// 获取帖子ID
	postID := c.Param("id")
	if postID == "" {
		logger.Errorf("LikePost postID is empty")
		common.RespondBadRequest(c, constants.MsgTargetIDRequired)
		return
	}

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		logger.Errorf("LikePost GetUserIDFromContext failed")
		common.RespondUnauthorized(c)
		return
	}

	// 构建请求
	req := &like.LikeRequest{
		UserId:     userID,
		TargetId:   postID,
		TargetType: constants.TargetTypePost,
	}

	// 调用点赞服务
	resp, err := handler.GetLikeClient().Like(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "LikePost", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "LikePost", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// UnlikePost 取消点赞帖子
func UnlikePost(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	// 获取帖子ID
	postID := c.Param("id")
	if postID == "" {
		logger.Errorf("UnlikePost postID is empty")
		common.RespondBadRequest(c, constants.MsgTargetIDRequired)
		return
	}

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		logger.Errorf("UnlikePost GetUserIDFromContext failed")
		common.RespondUnauthorized(c)
		return
	}

	// 构建请求
	req := &like.UnlikeRequest{
		UserId:     userID,
		TargetId:   postID,
		TargetType: constants.TargetTypePost,
	}

	// 调用取消点赞服务
	resp, err := handler.GetLikeClient().Unlike(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "UnlikePost", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "UnlikePost", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// GetPostLikeCount 获取帖子点赞数量
func GetPostLikeCount(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	// 获取帖子ID
	postID := c.Param("id")
	if postID == "" {
		logger.Errorf("GetPostLikeCount postID is empty")
		common.RespondBadRequest(c, constants.MsgTargetIDRequired)
		return
	}

	// 构建请求
	req := &like.GetLikeCountRequest{
		TargetId:   postID,
		TargetType: constants.TargetTypePost,
	}

	// 调用点赞服务
	resp, err := handler.GetLikeClient().GetLikeCount(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "GetPostLikeCount", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetPostLikeCount", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// CheckPostLikeStatus 检查帖子点赞状态
func CheckPostLikeStatus(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	// 获取帖子ID
	postID := c.Param("id")
	if postID == "" {
		logger.Errorf("CheckPostLikeStatus postID is empty")
		common.RespondBadRequest(c, constants.MsgTargetIDRequired)
		return
	}

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		logger.Errorf("CheckPostLikeStatus GetUserIDFromContext failed")
		common.RespondUnauthorized(c)
		return
	}

	// 构建请求
	req := &like.LikeRequest{
		UserId:     userID,
		TargetId:   postID,
		TargetType: constants.TargetTypePost,
	}

	// 调用点赞服务
	resp, err := handler.GetLikeClient().IsLiked(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "CheckPostLikeStatus", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "CheckPostLikeStatus", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// ==================== 评论点赞实现函数 ====================

// LikeComment 点赞评论
func LikeComment(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	// 获取评论ID
	commentID := c.Param("id")
	if commentID == "" {
		logger.Errorf("LikeComment commentID is empty")
		common.RespondBadRequest(c, constants.MsgTargetIDRequired)
		return
	}

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		logger.Errorf("LikeComment GetUserIDFromContext failed")
		common.RespondUnauthorized(c)
		return
	}

	// 构建请求
	req := &like.LikeRequest{
		UserId:     userID,
		TargetId:   commentID,
		TargetType: constants.TargetTypeComment,
	}

	// 调用点赞服务
	resp, err := handler.GetLikeClient().Like(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "LikeComment", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "LikeComment", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// UnlikeComment 取消点赞评论
func UnlikeComment(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	// 获取评论ID
	commentID := c.Param("id")
	if commentID == "" {
		logger.Errorf("UnlikeComment commentID is empty")
		common.RespondBadRequest(c, constants.MsgTargetIDRequired)
		return
	}

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		logger.Errorf("UnlikeComment GetUserIDFromContext failed")
		common.RespondUnauthorized(c)
		return
	}

	// 构建请求
	req := &like.UnlikeRequest{
		UserId:     userID,
		TargetId:   commentID,
		TargetType: constants.TargetTypeComment,
	}

	// 调用取消点赞服务
	resp, err := handler.GetLikeClient().Unlike(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "UnlikeComment", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "UnlikeComment", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// GetCommentLikeCount 获取评论点赞数量
func GetCommentLikeCount(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	// 获取评论ID
	commentID := c.Param("id")
	if commentID == "" {
		logger.Errorf("GetCommentLikeCount commentID is empty")
		common.RespondBadRequest(c, constants.MsgTargetIDRequired)
		return
	}

	// 构建请求
	req := &like.GetLikeCountRequest{
		TargetId:   commentID,
		TargetType: constants.TargetTypeComment,
	}

	// 调用点赞服务
	resp, err := handler.GetLikeClient().GetLikeCount(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "GetCommentLikeCount", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetCommentLikeCount", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// CheckCommentLikeStatus 检查评论点赞状态
func CheckCommentLikeStatus(ctx context.Context, c *app.RequestContext) {
	traceId, _ := c.Get(constants.TraceIdKey)
	logger := log.GetLogger().WithField(constants.TraceIdKey, traceId)
	// 获取评论ID
	commentID := c.Param("id")
	if commentID == "" {
		logger.Errorf("CheckCommentLikeStatus commentID is empty")
		common.RespondBadRequest(c, constants.MsgTargetIDRequired)
		return
	}

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		logger.Errorf("CheckCommentLikeStatus GetUserIDFromContext failed")
		common.RespondUnauthorized(c)
		return
	}

	// 构建请求
	req := &like.LikeRequest{
		UserId:     userID,
		TargetId:   commentID,
		TargetType: constants.TargetTypeComment,
	}

	// 调用点赞服务
	resp, err := handler.GetLikeClient().IsLiked(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "CheckCommentLikeStatus", traceId.(string))
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "CheckCommentLikeStatus", traceId.(string), resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// ==================== 通用点赞实现函数 ====================

// GetLikeList 获取用户点赞列表
func GetLikeList(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] GetLikeList request started", traceId)

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
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

	// 调用点赞服务
	resp, err := handler.GetLikeClient().GetLikeList(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "GetLikeList", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetLikeList", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// IsLiked 检查是否已点赞
func IsLiked(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] IsLiked request started", traceId)

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 获取查询参数
	targetID := c.Query("target_id")
	targetType := c.Query("target_type")

	// 验证参数
	if err := ValidateTargetParams(targetID, targetType); err != nil {
		common.RespondBadRequest(c, err.Error())
		return
	}

	// 构建请求
	req := &like.LikeRequest{
		UserId:     userID,
		TargetId:   targetID,
		TargetType: targetType,
	}

	// 调用点赞服务
	resp, err := handler.GetLikeClient().IsLiked(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "IsLiked", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "IsLiked", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// GetLikeCount 获取点赞数量
func GetLikeCount(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] GetLikeCount request started", traceId)

	// 获取查询参数
	targetID := c.Query("target_id")
	targetType := c.Query("target_type")

	// 验证参数
	if err := ValidateTargetParams(targetID, targetType); err != nil {
		common.RespondBadRequest(c, err.Error())
		return
	}

	// 构建请求
	req := &like.GetLikeCountRequest{
		TargetId:   targetID,
		TargetType: targetType,
	}

	// 调用点赞服务
	resp, err := handler.GetLikeClient().GetLikeCount(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "GetLikeCount", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetLikeCount", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// GetLikeUsers 获取点赞用户列表
func GetLikeUsers(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] GetLikeUsers request started", traceId)

	// 获取查询参数
	targetID := c.Query("target_id")
	targetType := c.Query("target_type")

	// 验证参数
	if err := ValidateTargetParams(targetID, targetType); err != nil {
		common.RespondBadRequest(c, err.Error())
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

	// 调用点赞服务
	resp, err := handler.GetLikeClient().GetLikeUsers(ctx, req)
	if err != nil {
		common.HandleRpcError(c, "GetLikeUsers", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "GetLikeUsers", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// ==================== 兼容性实现函数 ====================

// Like 兼容旧的点赞接口
func Like(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] Like request started", traceId)

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 解析请求参数
	var req like.LikeRequest
	if err := c.BindAndValidate(&req); err != nil {
		common.RespondBadRequest(c, constants.MsgParamError+": "+err.Error())
		return
	}

	// 设置用户ID
	req.UserId = userID

	// 验证目标类型
	if err := ValidateTargetParams(req.TargetId, req.TargetType); err != nil {
		common.RespondBadRequest(c, err.Error())
		return
	}

	// 调用点赞服务
	resp, err := handler.GetLikeClient().Like(ctx, &req)
	if err != nil {
		common.HandleRpcError(c, "Like", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "Like", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// Unlike 兼容旧的取消点赞接口
func Unlike(ctx context.Context, c *app.RequestContext) {
	// 获取trace ID
	traceId := c.GetString("trace_id")
	log.GetLogger().Infof("[%s] Unlike request started", traceId)

	// 获取用户ID
	userID, exists := common.GetUserIDFromContext(c)
	if !exists {
		common.RespondUnauthorized(c)
		return
	}

	// 解析请求参数
	var req like.UnlikeRequest
	if err := c.BindAndValidate(&req); err != nil {
		log.GetLogger().Errorf("[%s] Bind request error: %v", traceId, err)
		common.RespondBadRequest(c, constants.MsgParamError+": "+err.Error())
		return
	}

	// 设置用户ID
	req.UserId = userID

	// 验证目标类型
	if err := ValidateTargetParams(req.TargetId, req.TargetType); err != nil {
		common.RespondBadRequest(c, err.Error())
		return
	}

	// 调用点赞服务
	resp, err := handler.GetLikeClient().Unlike(ctx, &req)
	if err != nil {
		common.HandleRpcError(c, "Unlike", traceId)
		return
	}
	if resp.Code != constants.SuccessCode {
		common.HandleServiceError(c, "Unlike", traceId, resp.Code, resp.Message)
		return
	}
	common.RespondWithSuccess(c, resp)
}
