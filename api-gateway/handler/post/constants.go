package post

// HTTP状态码常量
const (
	HTTPStatusOK                  = 200
	HTTPStatusBadRequest          = 400
	HTTPStatusUnauthorized        = 401
	HTTPStatusForbidden           = 403
	HTTPStatusNotFound            = 404
	HTTPStatusInternalServerError = 500
)

// 响应消息常量
const (
	MsgSuccess            = "success"
	MsgUnauthorized       = "未授权"
	MsgParamError         = "参数错误"
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

// 默认分页参数
const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 50
)

// 排序类型常量
const (
	SortTypeLatest  = "latest"
	SortTypeHot     = "hot"
	SortTypeScore   = "score"
	SortTypeDaily   = "daily"
	SortTypeWeekly  = "weekly"
	SortTypeMonthly = "monthly"
)

// 排行榜类型常量
const (
	RankTypeDailyHigh     = "daily_high"
	RankTypeDailyLow      = "daily_low"
	RankTypeWeeklyBest    = "weekly_best"
	RankTypeControversial = "controversial"
)

// 话题排序类型常量
const (
	TopicSortTypeHot         = "hot"
	TopicSortTypeLatest      = "latest"
	TopicSortTypeParticipant = "participant"
)

// 上下文键常量
const (
	UserIDKey = "user_id"
)

// 查询参数键常量
const (
	ParamPage        = "page"
	ParamSize        = "size"
	ParamPageSize    = "page_size"
	ParamUserID      = "user_id"
	ParamTopicID     = "topic_id"
	ParamCategory    = "category"
	ParamTag         = "tag"
	ParamSortType    = "sort_type"
	ParamIsAnonymous = "is_anonymous"
	ParamRankType    = "rank_type"
	ParamDate        = "date"
	ParamKeyword     = "keyword"
	ParamLimit       = "limit"
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
