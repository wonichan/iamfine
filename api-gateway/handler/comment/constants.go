package comment

import "hupu/api-gateway/handler/common"

// 引用共享常量
const (
	// HTTP状态码
	HTTPStatusOK                  = common.HTTPStatusOK
	HTTPStatusBadRequest          = common.HTTPStatusBadRequest
	HTTPStatusUnauthorized        = common.HTTPStatusUnauthorized
	HTTPStatusForbidden           = common.HTTPStatusForbidden
	HTTPStatusNotFound            = common.HTTPStatusNotFound
	HTTPStatusInternalServerError = common.HTTPStatusInternalServerError

	// 业务状态码
CodeSuccess = common.CodeSuccess
CodeError   = common.CodeError

// 通用消息
MsgSuccess      = common.MsgSuccess
MsgUnauthorized = common.MsgUnauthorized
MsgParamError   = common.MsgParamError

	// 分页参数
	DefaultPage     = common.DefaultPage
	DefaultPageSize = common.DefaultPageSize
	MaxPageSize     = common.MaxPageSize

	// 上下文键
	ContextKeyUserID = common.ContextKeyUserID

	// 排序类型
	SortTypeLatest = common.SortTypeLatest
	SortTypeHot    = common.SortTypeHot
)

// 排序类型常量（评论特有）
const (
	SortTypeOldest = "oldest"
)

// 评论服务特有的错误消息常量
const (
	MsgCommentIDEmpty        = "评论ID不能为空"
	MsgPostIDEmpty           = "帖子ID不能为空"
	MsgUserIDEmpty           = "用户ID不能为空"
	MsgCreateCommentFailed   = "创建评论失败"
	MsgGetCommentListFailed  = "获取评论列表失败"
	MsgGetCommentFailed      = "获取评论详情失败"
	MsgDeleteCommentFailed   = "删除评论失败"
	MsgLikeCommentFailed     = "点赞评论失败"
	MsgUnlikeCommentFailed   = "取消点赞评论失败"
	MsgGetUserCommentsFailed = "获取用户评论列表失败"
	MsgNoPermissionDelete    = "无权限删除此评论"
)