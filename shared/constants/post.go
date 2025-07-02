package constants

// 帖子服务特有的错误消息常量
const (
	MsgPostIDEmpty        = "帖子ID不能为空"
	MsgTopicIDEmpty       = "话题ID不能为空"
	MsgKeywordEmpty       = "搜索关键词不能为空"
	MsgCreatePostFailed   = "创建帖子失败"
	MsgGetPostFailed      = "获取帖子失败"
	MsgUpdatePostFailed   = "更新帖子失败"
	MsgDeletePostFailed   = "删除帖子失败"
	MsgGetPostListFailed  = "获取帖子列表失败"
	MsgSearchPostsFailed  = "搜索帖子失败"
	MsgCreateTopicFailed  = "创建话题失败"
	MsgGetTopicFailed     = "获取话题失败"
	MsgGetTopicListFailed = "获取话题列表失败"
	MsgSearchTopicsFailed = "搜索话题失败"
	MsgShareTopicFailed   = "分享话题失败"
	MsgCollectPostFailed  = "收藏帖子失败"
	MsgUncollectFailed    = "取消收藏失败"
	MsgGetCollectedFailed = "获取收藏列表失败"
	MsgRatePostFailed     = "评分失败"
	MsgGetRatingFailed    = "获取评分失败"
	MsgUpdateRatingFailed = "更新评分失败"
	MsgDeleteRatingFailed = "删除评分失败"
	MsgGetRankFailed      = "获取排行榜失败"
)

type PostTopic int32

const (
	PostTopicUnknown PostTopic = iota
	PostTopicWarm
	PostTopicFunny
)

type PostCategory int32

const (
	PostCategoryUnknown PostCategory = iota
	PostCategoryRecommend
	PostCategoryHot
	PostCategoryDailyShare
	PostCategoryTodayHighScore
)
