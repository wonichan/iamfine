package constants

const (
	LikeMethodName            = "Like"
	UnlikeMethodName          = "Unlike"
	GetLikeListMethodName     = "GetLikeList"
	CheckLikeStatusMethodName = "CheckLikeStatus"
	GetLikeCountMethodName    = "GetLikeCount"
	GetLikeUsersMethodName    = "GetLikeUsers"
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
