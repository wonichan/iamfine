package constants

var FollowServiceErrorMsg = map[string]string{
	FollowMethodName:            MsgFollowFailed,
	UnFollowMethodName:          MsgUnfollowFailed,
	GetFollowListMethodName:     MsgGetFollowListFailed,
	GetFollowCountMethodName:    MsgGetFollowCountFailed,
	GetFollowerCountMethodName:  MsgGetFollowerCountFailed,
	GetMutualFollowsMethodName:  MsgGetMutualFollowsFailed,
	CheckFollowStatusMethodName: MsgCheckFollowStatusFailed,
	IsFollowingMethodName:       MsgCheckFollowStatusFailed,
}

var CommentServiceErrorMsg = map[string]string{
	CreateCommentMethodName:   MsgCreateCommentFailed,
	GetCommentListMethodName:  MsgGetCommentListFailed,
	GetCommentMethodName:      MsgGetCommentFailed,
	DeleteCommentMethodName:   MsgDeleteCommentFailed,
	GetUserCommentsMethodName: MsgGetUserCommentsFailed,
}

var LikeServiceErrorMsg = map[string]string{
	LikeMethodName:            MsgLikeFailed,
	UnlikeMethodName:          MsgUnlikeFailed,
	GetLikeListMethodName:     MsgGetLikeListFailed,
	CheckLikeStatusMethodName: MsgCheckLikeStatusFailed,
	GetLikeCountMethodName:    MsgGetLikeCountFailed,
	GetLikeUsersMethodName:    MsgGetLikeUsersFailed,
}
