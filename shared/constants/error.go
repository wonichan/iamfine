package constants

const SuccessCode = 0

const (
	InternalErrCode = 500
)
const (
	UserPasswdErrCode          = 100101
	UserExistsErrCode          = 100102
	UserNotExistsErrCode       = 100103
	UserFollowUserErrCode      = 100104
	UserUnfollowUserErrCode    = 100105
	UserLoginErrCode           = 100106
	UserGerUserErrCode         = 100107
	UserUpdateUserErrCode      = 100108
	UserGetFollowersErrCode    = 100109
	UserGetFollowingErrCode    = 100110
	UserCreateAnonymousErrCode = 100111
	UserGetAnonymousErrCode    = 100112
	UserUpdateAnonymousErrCode = 100113
	UserGetUserStatusErrCode   = 100114

	PostCreateErrCode       = 100201
	PostGetErrCode          = 100202
	PostUpdateErrCode       = 100203
	PostDeleteErrCode       = 100204
	PostNotFoundErrCode     = 100205
	PostListErrCode         = 100206
	PostCollectErrCode      = 100207
	PostUncollectErrCode    = 100208
	PostGetCollectedErrCode = 100209

	TopicCreateErrCode = 100301
	TopicGetErrCode    = 100302
	TopicListErrCode   = 100303

	RatePostErrCode          = 100401
	RatePostNotFoundErrCode  = 100402
	RateGetRatingRankErrCode = 100403

	ParamErrCode   = 500101
	InsertErrCode  = 500102
	UpdateErrCode  = 500103
	DeleteErrCode  = 500104
	DecryptErrCode = 500105
	EncryptErrCode = 500106
)

func ErrCodeToMsg(code int32) string {
	if msg, ok := errCodeToMsg[code]; ok {
		return msg
	}
	return MsgServerError
}
