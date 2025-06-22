package constants

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
