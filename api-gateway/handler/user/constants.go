package user

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
	MsgSuccess            = common.MsgSuccess
	MsgParamError         = common.MsgParamError
	MsgUnauthorized       = common.MsgUnauthorized
	MsgForbidden          = common.MsgForbidden
	MsgRequestFormatError = common.MsgRequestFormatError

	// 分页参数
	DefaultPage     = common.DefaultPage
	DefaultPageSize = common.DefaultPageSize

	// 上下文键
	ContextKeyUserID = common.ContextKeyUserID
)

// 用户服务特有的错误消息常量
const (
	MsgUserIDEmpty             = "用户ID不能为空"
	MsgTargetUserIDEmpty       = "目标用户ID不能为空"
	MsgProfileIDEmpty          = "profile_id不能为空"
	MsgCannotFollowSelf        = "不能关注自己"
	MsgNoPermissionUpdateOther = "无权限更新其他用户信息"
	MsgWxLoginFailed           = "微信登录失败"
	MsgRegisterFailed          = "注册失败"
	MsgLoginFailed             = "登录失败"
	MsgGetUserInfoFailed       = "获取用户信息失败"
	MsgUpdateUserInfoFailed    = "更新用户信息失败"
	MsgGetUserStatsFailed      = "获取用户统计信息失败"
	MsgGetUnreadCountFailed    = "获取未读消息数失败"
	MsgFollowUserFailed        = "关注用户失败"
	MsgUnfollowUserFailed      = "取消关注用户失败"
	MsgGetFollowersFailed      = "获取粉丝列表失败"
	MsgGetFollowingFailed      = "获取关注列表失败"
	MsgCreateAnonymousFailed   = "创建匿名马甲失败"
	MsgGetAnonymousListFailed  = "获取匿名马甲列表失败"
	MsgUpdateAnonymousFailed   = "更新匿名马甲失败"
)
