package user

// HTTP状态码常量
const (
	HTTPStatusOK                  = 200
	HTTPStatusBadRequest          = 400
	HTTPStatusUnauthorized        = 401
	HTTPStatusForbidden           = 403
	HTTPStatusNotFound            = 404
	HTTPStatusInternalServerError = 500
)

// 业务状态码常量
const (
	CodeSuccess = 200
	CodeError   = 500
)

// 错误消息常量
const (
	MsgSuccess                 = "success"
	MsgParamError              = "参数错误"
	MsgUnauthorized            = "未授权"
	MsgForbidden               = "禁止访问"
	MsgUserIDEmpty             = "用户ID不能为空"
	MsgTargetUserIDEmpty       = "目标用户ID不能为空"
	MsgProfileIDEmpty          = "profile_id不能为空"
	MsgCannotFollowSelf        = "不能关注自己"
	MsgRequestFormatError      = "请求参数格式错误"
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

// 默认分页参数
const (
	DefaultPage     = 1
	DefaultPageSize = 10
)

// 上下文键常量
const (
	ContextKeyUserID = "user_id"
)
