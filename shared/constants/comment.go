package constants

// 评论服务特有的错误消息常量
const (
	MsgCommentIDEmpty        = "评论ID不能为空"
	MsgCreateCommentFailed   = "创建评论失败"
	MsgGetCommentListFailed  = "获取评论列表失败"
	MsgGetCommentFailed      = "获取评论详情失败"
	MsgDeleteCommentFailed   = "删除评论失败"
	MsgLikeCommentFailed     = "点赞评论失败"
	MsgUnlikeCommentFailed   = "取消点赞评论失败"
	MsgGetUserCommentsFailed = "获取用户评论列表失败"
	MsgNoPermissionDelete    = "无权限删除此评论"
)

// 排序类型常量（评论特有）
const (
	SortTypeOldest = "oldest"
)
