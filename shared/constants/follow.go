package constants

import "hupu/api-gateway/handler/common"

// 引用共享常量
const (
	// HTTP Status Codes
	HTTPStatusOK                  = common.HTTPStatusOK
	HTTPStatusBadRequest          = common.HTTPStatusBadRequest
	HTTPStatusUnauthorized        = common.HTTPStatusUnauthorized
	HTTPStatusForbidden           = common.HTTPStatusForbidden
	HTTPStatusNotFound            = common.HTTPStatusNotFound
	HTTPStatusInternalServerError = common.HTTPStatusInternalServerError

	// Default pagination values
	DefaultPage     = common.DefaultPage
	DefaultPageSize = common.DefaultPageSize
	MaxPageSize     = common.MaxPageSize

	// Follow status
	FollowStatusNotFollowing = common.FollowStatusNotFollowing
	FollowStatusFollowing    = common.FollowStatusFollowing
	FollowStatusMutual       = common.FollowStatusMutual

	// 通用消息
	MsgUnauthorized = common.MsgUnauthorized
	MsgParamError   = common.MsgParamError
	MsgSuccess      = common.MsgSuccess
)

// 关注服务特有的错误消息常量
const (
	MsgUserIDRequired         = "用户ID不能为空"
	MsgTargetUserIDRequired   = "目标用户ID不能为空"
	MsgFollowFailed           = "关注失败"
	MsgUnfollowFailed         = "取消关注失败"
	MsgGetFollowListFailed    = "获取关注列表失败"
	MsgGetFollowerListFailed  = "获取粉丝列表失败"
	MsgCheckFollowStatusFailed = "检查关注状态失败"
	MsgGetFollowCountFailed   = "获取关注数量失败"
	MsgGetFollowerCountFailed = "获取粉丝数量失败"
	MsgGetMutualFollowsFailed = "获取共同关注失败"
)