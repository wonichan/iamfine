package like

import "hupu/api-gateway/handler/common"

// 引用共享常量
const (
	// HTTP状态码
	HTTPStatusOK                  = common.HTTPStatusOK
	HTTPStatusBadRequest          = common.HTTPStatusBadRequest
	HTTPStatusUnauthorized        = common.HTTPStatusUnauthorized
	HTTPStatusInternalServerError = common.HTTPStatusInternalServerError

	// 业务状态码
CodeSuccess     = common.CodeSuccess
CodeError       = common.CodeError
CodeServerError = common.CodeServerError

	// 通用消息
	MsgUnauthorized = common.MsgUnauthorized
	MsgParamError   = common.MsgParamError

	// 分页参数
	DefaultPage     = common.DefaultPage
	DefaultPageSize = common.DefaultPageSize
	MaxPageSize     = common.MaxPageSize

	// 上下文键
	UserIDKey = common.UserIDKey
)

// 点赞服务特有的错误消息常量
const (
	MsgTargetIDRequired      = "目标ID不能为空"
	MsgTargetTypeRequired    = "目标类型不能为空"
	MsgLikeFailed            = "点赞失败"
	MsgUnlikeFailed          = "取消点赞失败"
	MsgGetLikeListFailed     = "获取点赞列表失败"
	MsgCheckLikeStatusFailed = "检查点赞状态失败"
	MsgGetLikeCountFailed    = "获取点赞数量失败"
	MsgGetLikeUsersFailed    = "获取点赞用户列表失败"
)

// 目标类型常量
const (
	TargetTypePost    = "post"
	TargetTypeComment = "comment"
)