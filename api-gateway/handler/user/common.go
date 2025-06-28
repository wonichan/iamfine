package user

import (
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/shared/log"
)

// ResponseData 通用响应结构
type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应
func ErrorResponse(c *app.RequestContext, httpCode int, code int, message string) {
	c.JSON(httpCode, ResponseData{
		Code:    code,
		Message: message,
	})
}

// SuccessResponse 成功响应
func SuccessResponse(c *app.RequestContext, data interface{}) {
	c.JSON(HTTPStatusOK, data)
}

// GetUserIDFromContext 从上下文获取用户ID
func GetUserIDFromContext(c *app.RequestContext) (string, bool) {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// ValidateUserIDParam 验证用户ID参数
func ValidateUserIDParam(c *app.RequestContext, paramName string) (string, bool) {
	userID := c.Param(paramName)
	if userID == "" {
		ErrorResponse(c, HTTPStatusBadRequest, CodeError, MsgUserIDEmpty)
		return "", false
	}
	return userID, true
}

// ValidateTargetUserIDParam 验证目标用户ID参数
func ValidateTargetUserIDParam(c *app.RequestContext, paramName string) (string, bool) {
	targetUserID := c.Param(paramName)
	if targetUserID == "" {
		ErrorResponse(c, HTTPStatusBadRequest, CodeError, MsgTargetUserIDEmpty)
		return "", false
	}
	return targetUserID, true
}

// ValidateProfileIDParam 验证profile ID参数
func ValidateProfileIDParam(c *app.RequestContext, paramName string) (string, bool) {
	profileID := c.Param(paramName)
	if profileID == "" {
		ErrorResponse(c, HTTPStatusBadRequest, CodeError, MsgProfileIDEmpty)
		return "", false
	}
	return profileID, true
}

// CheckUserPermission 检查用户权限（只能操作自己的信息）
func CheckUserPermission(c *app.RequestContext, currentUserID, targetUserID string) bool {
	if currentUserID != targetUserID {
		ErrorResponse(c, HTTPStatusForbidden, CodeError, MsgNoPermissionUpdateOther)
		return false
	}
	return true
}

// CheckSelfFollow 检查是否关注自己
func CheckSelfFollow(c *app.RequestContext, userID, targetUserID string) bool {
	if userID == targetUserID {
		ErrorResponse(c, HTTPStatusBadRequest, CodeError, MsgCannotFollowSelf)
		return false
	}
	return true
}

// ParsePaginationParams 解析分页参数
func ParsePaginationParams(c *app.RequestContext) (int32, int32) {
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = DefaultPage
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = DefaultPageSize
	}

	return int32(page), int32(pageSize)
}

// LogError 记录错误日志
func LogError(operation string, err error) {
	log.GetLogger().Errorf("%s error: %v", operation, err)
}

// HandleServiceError 处理服务调用错误
func HandleServiceError(c *app.RequestContext, operation string, err error, message string) {
	LogError(operation, err)
	ErrorResponse(c, HTTPStatusInternalServerError, CodeError, message)
}

// RequireAuth 需要认证的接口通用检查
func RequireAuth(c *app.RequestContext) (string, bool) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		ErrorResponse(c, HTTPStatusUnauthorized, CodeError, MsgUnauthorized)
		return "", false
	}
	return userID, true
}
