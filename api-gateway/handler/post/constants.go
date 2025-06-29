package post

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

	// 通用消息
	MsgSuccess      = common.MsgSuccess
	MsgUnauthorized = common.MsgUnauthorized
	MsgParamError   = common.MsgParamError

	// 分页参数
	DefaultPage     = common.DefaultPage
	DefaultPageSize = common.DefaultPageSize
	MaxPageSize     = common.MaxPageSize
)

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
	UserIDKey = common.UserIDKey
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
