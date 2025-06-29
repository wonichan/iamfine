package constants

const (
	FollowMethodName            = "Follow"
	UnFollowMethodName          = "Unfollow"
	GetFollowListMethodName     = "GetFollowList"
	GetFollowCountMethodName    = "GetFollowCount"
	GetFollowerCountMethodName  = "GetFollowerCount"
	GetMutualFollowsMethodName  = "GetMutualFollows"
	CheckFollowStatusMethodName = "CheckFollowStatus"
	IsFollowingMethodName       = "IsFollowing"
)

// 关注服务特有的错误消息常量
const (
	MsgUserIDRequired          = "用户ID不能为空"
	MsgTargetUserIDRequired    = "目标用户ID不能为空"
	MsgFollowFailed            = "关注失败"
	MsgUnfollowFailed          = "取消关注失败"
	MsgGetFollowListFailed     = "获取关注列表失败"
	MsgGetFollowerListFailed   = "获取粉丝列表失败"
	MsgCheckFollowStatusFailed = "检查关注状态失败"
	MsgGetFollowCountFailed    = "获取关注数量失败"
	MsgGetFollowerCountFailed  = "获取粉丝数量失败"
	MsgGetMutualFollowsFailed  = "获取共同关注失败"
)
