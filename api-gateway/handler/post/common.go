package post

import (
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/shared/log"
)

// 公共响应结构
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 获取用户ID的公共方法
func GetUserID(c *app.RequestContext) (string, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// 返回未授权错误的公共方法
func RespondUnauthorized(c *app.RequestContext) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Code:    HTTPStatusUnauthorized,
		Message: MsgUnauthorized,
	})
}

// 返回参数错误的公共方法
func RespondBadRequest(c *app.RequestContext, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    HTTPStatusBadRequest,
		Message: message,
	})
}

// 返回内部服务器错误的公共方法
func RespondInternalError(c *app.RequestContext, message string, err error) {
	log.GetLogger().Errorf("%s: %v", message, err)
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    HTTPStatusInternalServerError,
		Message: message,
	})
}

// 解析分页参数的公共方法
func ParsePaginationParams(c *app.RequestContext) (page, pageSize int64) {
	pageStr := c.Query(ParamPage)
	pageSizeStr := c.Query(ParamPageSize)

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = DefaultPage
	}

	pageSize, err = strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return page, pageSize
}

// 解析分页参数（返回int32）的公共方法
func ParsePaginationParamsInt32(c *app.RequestContext) (page, pageSize int32) {
	p, ps := ParsePaginationParams(c)
	return int32(p), int32(ps)
}

// 获取路径参数的公共方法
func GetPathParam(c *app.RequestContext, key string) string {
	return c.Param(key)
}

// 验证必需路径参数的公共方法
func ValidateRequiredPathParam(c *app.RequestContext, key, paramName string) (string, bool) {
	value := c.Param(key)
	if value == "" {
		RespondBadRequest(c, paramName+"不能为空")
		return "", false
	}
	return value, true
}

// 解析可选字符串参数的公共方法
func ParseOptionalStringParam(c *app.RequestContext, key string) *string {
	value := c.Query(key)
	if value == "" {
		return nil
	}
	return &value
}

// 解析可选布尔参数的公共方法
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

// 解析可选整数参数的公共方法
func ParseOptionalIntParam(c *app.RequestContext, key string) *int32 {
	valueStr := c.Query(key)
	if valueStr == "" {
		return nil
	}
	if value, err := strconv.ParseInt(valueStr, 10, 32); err == nil {
		intValue := int32(value)
		return &intValue
	}
	return nil
}

// 绑定和验证请求参数的公共方法
func BindAndValidateRequest(c *app.RequestContext, req interface{}) bool {
	if err := c.BindAndValidate(req); err != nil {
		RespondBadRequest(c, MsgParamError+": "+err.Error())
		return false
	}
	return true
}
