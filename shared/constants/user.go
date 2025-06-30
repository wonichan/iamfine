package constants

const (
	UserStatusUnknown = 0
	UserStatusNormal  = iota
	UserStatusLocked
	UserStatusDeleted
	UserStatusForbidden
	UserStatusLogin
)

// 用户服务特有的错误消息常量
const (
	MsgPasswordError           = "密码错误"
	MsgEncryptError            = "加密错误"
	MsgDecryptError            = "解密错误"
	MsgUserExists              = "用户已存在"
	MsgUserNotExists           = "用户不存在"
	MsgUserStatusError         = "用户状态异常"
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
	MsgGetAnonymousFailed      = "获取匿名马甲失败"
	MsgUpdateAnonymousFailed   = "更新匿名马甲失败"
)

var errCodeToMsg = map[int32]string{
	InternalErrCode:      MsgServerError,
	UserPasswdErrCode:    MsgPasswordError,
	UserExistsErrCode:    MsgUserExists,
	UserNotExistsErrCode: MsgUserNotExists,
	ParamErrCode:         MsgParamError,
	InsertErrCode:        InsertError,
	UpdateErrCode:        UpdateError,
	DeleteErrCode:        DeleteError,
	DecryptErrCode:       MsgDecryptError,
	EncryptErrCode:       MsgEncryptError,
}
