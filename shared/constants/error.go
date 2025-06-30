package constants

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

	ParamErrCode   = 500101
	InsertErrCode  = 500102
	UpdateErrCode  = 500103
	DeleteErrCode  = 500104
	DecryptErrCode = 500105
	EncryptErrCode = 500106
)

func ErrMsgToCode(msg string) int32 {
	if code, ok := errMsgToCode[msg]; ok {
		return code
	}
	return InternalErrCode
}
