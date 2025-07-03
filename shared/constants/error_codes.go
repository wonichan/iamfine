package constants

// 基础错误码
const (
	// 成功
	SuccessCode = 0

	// 通用错误码 (1000-1999)
	InternalServerErrorCode = 1000
	InvalidParameterCode    = 1001
	ValidationErrorCode     = 1002
	DatabaseErrorCode       = 1003
	CacheErrorCode          = 1004
	NetworkErrorCode        = 1005
	TimeoutErrorCode        = 1006
	PermissionDeniedCode    = 1007
	ResourceNotFoundCode    = 1008
	ResourceConflictCode    = 1009
	RateLimitExceededCode   = 1010
)

// 用户相关错误码 (2000-2999)
const (
	UserNotFoundCode           = 2000
	UserAlreadyExistsCode      = 2001
	UserInvalidCredentialsCode = 2002
	UserAccountDisabledCode    = 2003
	UserPasswordTooWeakCode    = 2004
	UserEmailInvalidCode       = 2005
	UserPhoneInvalidCode       = 2006
	UserTokenExpiredCode       = 2007
	UserTokenInvalidCode       = 2008
	UserFollowSelfCode         = 2009
	UserAlreadyFollowedCode    = 2010
	UserFollowFailCode         = 2011
	UserNotFollowedCode        = 2012
	UserProfileUpdateFailCode  = 2013
)

// 帖子相关错误码 (3000-3999)
const (
	PostNotFoundCode         = 3000
	PostCreateFailCode       = 3001
	PostUpdateFailCode       = 3002
	PostDeleteFailCode       = 3003
	PostTitleTooLongCode     = 3004
	PostContentTooLongCode   = 3005
	PostTitleEmptyCode       = 3006
	PostContentEmptyCode     = 3007
	PostAlreadyFavoritedCode = 3008
	PostNotFavoritedCode     = 3009
	PostFavoriteFailCode     = 3010
	PostUnfavoriteFailCode   = 3011
	PostRatingFailCode       = 3012
	PostViewCountFailCode    = 3013
	PostRatingRankFailCode   = 3014
)

// 话题相关错误码 (4000-4999)
const (
	TopicNotFoundCode      = 4000
	TopicCreateFailCode    = 4001
	TopicUpdateFailCode    = 4002
	TopicDeleteFailCode    = 4003
	TopicNameEmptyCode     = 4004
	TopicNameTooLongCode   = 4005
	TopicAlreadyExistsCode = 4006
)

// 评论相关错误码 (5000-5999)
const (
	CommentNotFoundCode         = 5000
	CommentCreateFailCode       = 5001
	CommentUpdateFailCode       = 5002
	CommentDeleteFailCode       = 5003
	CommentContentEmptyCode     = 5004
	CommentContentTooLongCode   = 5005
	CommentPermissionDeniedCode = 5006
)

// 点赞相关错误码 (6000-6999)
const (
	LikeNotFoundCode      = 6000
	LikeCreateFailCode    = 6001
	LikeDeleteFailCode    = 6002
	LikeAlreadyExistsCode = 6003
	LikeNotExistsCode     = 6004
	LikeSelfPostCode      = 6005
)

// 关注相关错误码 (7000-7999)
const (
	FollowNotFoundCode      = 7000
	FollowCreateFailCode    = 7001
	FollowDeleteFailCode    = 7002
	FollowAlreadyExistsCode = 7003
	FollowNotExistsCode     = 7004
	FollowSelfCode          = 7005
	FollowListFailCode      = 7006
	FollowerListFailCode    = 7007
)

// 通知相关错误码 (8000-8999)
const (
	NotificationNotFoundCode     = 8000
	NotificationCreateFailCode   = 8001
	NotificationUpdateFailCode   = 8002
	NotificationDeleteFailCode   = 8003
	NotificationMarkReadFailCode = 8004
	NotificationListFailCode     = 8005
)

// 文件上传相关错误码 (9000-9999)
const (
	FileUploadFailCode       = 9000
	FileNotFoundCode         = 9001
	FileDeleteFailCode       = 9002
	FileSizeTooLargeCode     = 9003
	FileTypeNotSupportedCode = 9004
	FileNameInvalidCode      = 9005
)

// ErrorCodeMessage 错误码对应的消息
var ErrorCodeMessage = map[int32]string{
	// 成功
	SuccessCode: "操作成功",

	// 通用错误
	InternalServerErrorCode: "内部服务器错误",
	InvalidParameterCode:    "参数无效",
	ValidationErrorCode:     "参数验证失败",
	DatabaseErrorCode:       "数据库操作失败",
	CacheErrorCode:          "缓存操作失败",
	NetworkErrorCode:        "网络错误",
	TimeoutErrorCode:        "操作超时",
	PermissionDeniedCode:    "权限不足",
	ResourceNotFoundCode:    "资源不存在",
	ResourceConflictCode:    "资源冲突",
	RateLimitExceededCode:   "请求频率超限",

	// 用户相关错误
	UserNotFoundCode:           "用户不存在",
	UserAlreadyExistsCode:      "用户已存在",
	UserInvalidCredentialsCode: "用户名或密码错误",
	UserAccountDisabledCode:    "账户已被禁用",
	UserPasswordTooWeakCode:    "密码强度不够",
	UserEmailInvalidCode:       "邮箱格式无效",
	UserPhoneInvalidCode:       "手机号格式无效",
	UserTokenExpiredCode:       "登录已过期",
	UserTokenInvalidCode:       "登录凭证无效",
	UserFollowSelfCode:         "不能关注自己",
	UserFollowFailCode:         "关注失败",
	UserAlreadyFollowedCode:    "已经关注过该用户",
	UserNotFollowedCode:        "未关注该用户",
	UserProfileUpdateFailCode:  "用户信息更新失败",

	// 帖子相关错误
	PostNotFoundCode:         "帖子不存在",
	PostCreateFailCode:       "帖子创建失败",
	PostUpdateFailCode:       "帖子更新失败",
	PostDeleteFailCode:       "帖子删除失败",
	PostTitleTooLongCode:     "帖子标题过长",
	PostContentTooLongCode:   "帖子内容过长",
	PostTitleEmptyCode:       "帖子标题不能为空",
	PostContentEmptyCode:     "帖子内容不能为空",
	PostAlreadyFavoritedCode: "已经收藏过该帖子",
	PostNotFavoritedCode:     "未收藏该帖子",
	PostFavoriteFailCode:     "收藏帖子失败",
	PostUnfavoriteFailCode:   "取消收藏失败",
	PostRatingFailCode:       "帖子评分失败",
	PostViewCountFailCode:    "浏览次数更新失败",
	PostRatingRankFailCode:   "帖子排名获取失败",

	// 话题相关错误
	TopicNotFoundCode:      "话题不存在",
	TopicCreateFailCode:    "话题创建失败",
	TopicUpdateFailCode:    "话题更新失败",
	TopicDeleteFailCode:    "话题删除失败",
	TopicNameEmptyCode:     "话题名称不能为空",
	TopicNameTooLongCode:   "话题名称过长",
	TopicAlreadyExistsCode: "话题已存在",

	// 评论相关错误
	CommentNotFoundCode:         "评论不存在",
	CommentCreateFailCode:       "评论创建失败",
	CommentUpdateFailCode:       "评论更新失败",
	CommentDeleteFailCode:       "评论删除失败",
	CommentContentEmptyCode:     "评论内容不能为空",
	CommentContentTooLongCode:   "评论内容过长",
	CommentPermissionDeniedCode: "无权限操作该评论",

	// 点赞相关错误
	LikeNotFoundCode:      "点赞记录不存在",
	LikeCreateFailCode:    "点赞失败",
	LikeDeleteFailCode:    "取消点赞失败",
	LikeAlreadyExistsCode: "已经点赞过",
	LikeNotExistsCode:     "未点赞",
	LikeSelfPostCode:      "不能给自己的帖子点赞",

	// 关注相关错误
	FollowNotFoundCode:      "关注记录不存在",
	FollowCreateFailCode:    "关注失败",
	FollowDeleteFailCode:    "取消关注失败",
	FollowAlreadyExistsCode: "已经关注过",
	FollowNotExistsCode:     "未关注",
	FollowSelfCode:          "不能关注自己",
	FollowListFailCode:      "获取关注列表失败",
	FollowerListFailCode:    "获取粉丝列表失败",

	// 通知相关错误
	NotificationNotFoundCode:     "通知不存在",
	NotificationCreateFailCode:   "通知创建失败",
	NotificationUpdateFailCode:   "通知更新失败",
	NotificationDeleteFailCode:   "通知删除失败",
	NotificationMarkReadFailCode: "标记已读失败",
	NotificationListFailCode:     "获取通知列表失败",

	// 文件上传相关错误
	FileUploadFailCode:       "文件上传失败",
	FileNotFoundCode:         "文件不存在",
	FileDeleteFailCode:       "文件删除失败",
	FileSizeTooLargeCode:     "文件大小超限",
	FileTypeNotSupportedCode: "文件类型不支持",
	FileNameInvalidCode:      "文件名无效",
}

// GetErrorMessage 根据错误码获取错误消息
func GetErrorMessage(code int32) string {
	if msg, exists := ErrorCodeMessage[code]; exists {
		return msg
	}
	return MsgServerError
}

// IsSuccessCode 判断是否为成功码
func IsSuccessCode(code int32) bool {
	return code == SuccessCode
}

// IsUserError 判断是否为用户相关错误
func IsUserError(code int32) bool {
	return code >= 2000 && code < 3000
}

// IsPostError 判断是否为帖子相关错误
func IsPostError(code int32) bool {
	return code >= 3000 && code < 4000
}

// IsTopicError 判断是否为话题相关错误
func IsTopicError(code int32) bool {
	return code >= 4000 && code < 5000
}

// IsCommentError 判断是否为评论相关错误
func IsCommentError(code int32) bool {
	return code >= 5000 && code < 6000
}

// IsLikeError 判断是否为点赞相关错误
func IsLikeError(code int32) bool {
	return code >= 6000 && code < 7000
}

// IsFollowError 判断是否为关注相关错误
func IsFollowError(code int32) bool {
	return code >= 7000 && code < 8000
}

// IsNotificationError 判断是否为通知相关错误
func IsNotificationError(code int32) bool {
	return code >= 8000 && code < 9000
}

// IsFileError 判断是否为文件相关错误
func IsFileError(code int32) bool {
	return code >= 9000 && code < 10000
}
