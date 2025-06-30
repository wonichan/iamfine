package constants

import "net/http"

const (
	TraceIdKey = "trace_id"
)

// 引用共享常量
const (
	// HTTP Status Codes
	HTTPStatusOK                  = http.StatusOK
	HTTPStatusBadRequest          = http.StatusBadRequest
	HTTPStatusUnauthorized        = http.StatusUnauthorized
	HTTPStatusForbidden           = http.StatusForbidden
	HTTPStatusNotFound            = http.StatusNotFound
	HTTPStatusInternalServerError = http.StatusInternalServerError

	// 分页参数
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100

	// 关注状态
	FollowStatusNotFollowing = 0 // 未关注
	FollowStatusFollowing    = 1 // 已关注
	FollowStatusMutual       = 2 // 相互关注

	// 通用消息
	MsgSuccess            = "操作成功"
	MsgParamError         = "参数错误"
	MsgUnauthorized       = "未授权访问"
	MsgForbidden          = "禁止访问"
	MsgServerError        = "服务器内部错误"
	MsgRequestFormatError = "请求格式错误"
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

const (
	InsertError = "数据插入失败"
	UpdateError = "数据更新失败"
	DeleteError = "数据删除失败"
)
