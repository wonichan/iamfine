package common

import (
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/shared/constants"
	"hupu/shared/log"
)

// 业务状态码常量
const (
	// 业务状态码
	CodeSuccess     = 0
	CodeError       = 200500
	CodeServerError = 500
)

// 上下文键常量
const (
	// 上下文键
	UserIDKey        = "user_id"
	ContextKeyUserID = "user_id"
	UserRoleKey      = "role"
)

// 关注状态常量
const ()

// 通用响应结构
type ResponseData struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 错误响应结构
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 成功响应结构
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponseFunc 错误响应函数
func ErrorResponseFunc(c *app.RequestContext, httpCode int, code int32, message string) {
	c.JSON(httpCode, ResponseData{
		Code:    code,
		Message: message,
	})
}

// SuccessResponseFunc 成功响应函数
func SuccessResponseFunc(c *app.RequestContext, data interface{}) {
	c.JSON(constants.HTTPStatusOK, data)
}

// GetUserIDFromContext 从上下文获取用户ID
func GetUserIDFromContext(c *app.RequestContext) (string, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GetUserRoleFromContext 从上下文获取用户角色
func GetUserRoleFromContext(c *app.RequestContext) (string, bool) {
	role, exists := c.Get(UserRoleKey)
	if !exists {
		return "", false
	}
	return role.(string), true
}

// RequireAuth 需要认证的接口通用检查
func RequireAuth(c *app.RequestContext) (string, bool) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		ErrorResponseFunc(c, constants.HTTPStatusUnauthorized, CodeError, constants.MsgUnauthorized)
		return "", false
	}
	return userID, true
}

// ParsePaginationParams 解析分页参数
func ParsePaginationParams(c *app.RequestContext) (int32, int32) {
	pageStr := c.Query(constants.ParamPage)
	pageSizeStr := c.Query(constants.ParamPageSize)

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = constants.DefaultPage
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = constants.DefaultPageSize
	}

	if pageSize > constants.MaxPageSize {
		pageSize = constants.MaxPageSize
	}

	return int32(page), int32(pageSize)
}

// ParseOptionalStringParam 解析可选字符串参数
func ParseOptionalStringParam(c *app.RequestContext, key string) *string {
	value := c.Query(key)
	if value == "" {
		return nil
	}
	return &value
}

// ParseOptionalBoolParam 解析可选布尔参数
func ParseOptionalBoolParam(c *app.RequestContext, key string) *bool {
	valueStr := c.Query(key)
	if valueStr == "" {
		return nil
	}
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return &value
	}
	return nil
}

// ParseOptionalIntParam 解析可选整数参数
func ParseOptionalIntParam(c *app.RequestContext, key string) int32 {
	valueStr := c.Query(key)
	if valueStr == "" {
		return 0
	}
	if value, err := strconv.ParseInt(valueStr, 10, 32); err == nil {
		intValue := int32(value)
		return intValue
	}
	return 0
}

// LogError 记录错误日志
func LogError(operation, traceID string, err string) {
	logger := log.GetLogger()
	if traceID != "" {
		logger.WithField(constants.TraceIdKey, traceID)
	}
	logger.Errorf("%s error: %s", operation, err)
}

// HandleServiceError 处理服务调用错误
func HandleServiceError(c *app.RequestContext, operation, traceID string, code int32, errMsg string) {
	LogError(operation, traceID, errMsg)
	ErrorResponseFunc(c, constants.HTTPStatusInternalServerError, code, constants.GetErrorMessage(code))
}

func HandleRpcError(c *app.RequestContext, operation, traceID string) {
	LogError(operation, traceID, constants.MsgServerError)
	HandleServiceError(c, operation, traceID, constants.InternalErrCode, constants.MsgServerError)
}

// ValidateRequiredParam 验证必需参数
func ValidateRequiredParam(c *app.RequestContext, value, paramName string) bool {
	if value == "" {
		ErrorResponseFunc(c, constants.HTTPStatusBadRequest, CodeError, paramName+"不能为空")
		return false
	}
	return true
}

// ValidateRequiredPathParam 验证必需路径参数
func ValidateRequiredPathParam(c *app.RequestContext, key, paramName string) (string, bool) {
	value := c.Param(key)
	if !ValidateRequiredParam(c, value, paramName) {
		return "", false
	}
	return value, true
}

// ValidateUserIDParam 验证用户ID参数
func ValidateUserIDParam(c *app.RequestContext, paramName string) (string, bool) {
	userID := c.Param(paramName)
	if !ValidateRequiredParam(c, userID, "用户ID") {
		return "", false
	}
	return userID, true
}

// ValidateTargetUserIDParam 验证目标用户ID参数
func ValidateTargetUserIDParam(c *app.RequestContext, paramName string) (string, bool) {
	targetUserID := c.Param(paramName)
	if !ValidateRequiredParam(c, targetUserID, "目标用户ID") {
		return "", false
	}
	return targetUserID, true
}

// ValidateProfileIDParam 验证profile ID参数
func ValidateProfileIDParam(c *app.RequestContext, paramName string) (string, bool) {
	profileID := c.Param(paramName)
	if !ValidateRequiredParam(c, profileID, "Profile ID") {
		return "", false
	}
	return profileID, true
}

// ValidateCommentIDParam 验证评论ID参数
func ValidateCommentIDParam(c *app.RequestContext, paramName string) (string, bool) {
	commentID := c.Param(paramName)
	if !ValidateRequiredParam(c, commentID, "评论ID") {
		return "", false
	}
	return commentID, true
}

// ValidatePostIDParam 验证帖子ID参数（从查询参数获取）
func ValidatePostIDParam(c *app.RequestContext, paramName string) (string, bool) {
	postID := c.Query(paramName)
	if !ValidateRequiredParam(c, postID, "帖子ID") {
		return "", false
	}
	return postID, true
}

// ValidatePostIDPathParam 验证帖子ID参数（从路径参数获取）
func ValidatePostIDPathParam(c *app.RequestContext, paramName string) (string, bool) {
	postID := c.Param(paramName)
	if !ValidateRequiredParam(c, postID, "帖子ID") {
		return "", false
	}
	return postID, true
}

// RespondUnauthorized 返回未授权错误
func RespondUnauthorized(c *app.RequestContext) {
	ErrorResponseFunc(c, constants.HTTPStatusUnauthorized, CodeError, constants.MsgUnauthorized)
}

// RespondBadRequest 返回参数错误
func RespondBadRequest(c *app.RequestContext, message string) {
	ErrorResponseFunc(c, constants.HTTPStatusBadRequest, CodeError, message)
}

// RespondWithSuccess 统一成功响应（别名）
func RespondWithSuccess(c *app.RequestContext, data interface{}) {
	SuccessResponseFunc(c, data)
}

// GetUserID 获取用户ID（别名）
func GetUserID(c *app.RequestContext) (string, bool) {
	return GetUserIDFromContext(c)
}

// GetPathParam 获取路径参数
func GetPathParam(c *app.RequestContext, key string) string {
	return c.Param(key)
}

// ParsePaginationParamsInt64 解析分页参数（返回int64）
func ParsePaginationParamsInt64(c *app.RequestContext) (page, pageSize int64) {
	p, ps := ParsePaginationParams(c)
	return int64(p), int64(ps)
}

// ParseOptionalBoolParamValue 解析可选布尔参数（返回bool而非*bool）
func ParseOptionalBoolParamValue(c *app.RequestContext, key string) bool {
	boolPtr := ParseOptionalBoolParam(c, key)
	if boolPtr == nil {
		return false
	}
	return *boolPtr
}

// BindAndValidateRequest 绑定和验证请求参数
func BindAndValidateRequest(c *app.RequestContext, req interface{}) bool {
	if err := c.BindAndValidate(req); err != nil {
		RespondBadRequest(c, constants.MsgParamError+": "+err.Error())
		return false
	}
	return true
}

// ServiceCall 服务调用函数类型
type ServiceCall func() (any, error)

// CallService 统一服务调用处理
func CallService(c *app.RequestContext, serviceCall ServiceCall, operation, errorMsg string) {
	// 获取 trace ID
	traceID := c.GetString(constants.TraceIdKey)

	// 执行服务调用
	result, err := serviceCall()
	if err != nil {
		HandleServiceError(c, operation, traceID, constants.InternalErrCode, errorMsg)
		return
	}

	// 返回成功响应
	SuccessResponseFunc(c, result)
}
