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
