package follow

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// ==================== 关注管理Handler ====================

// FollowHandler 关注用户
func FollowHandler(ctx context.Context, c *app.RequestContext) {
	Follow(ctx, c)
}

// UnfollowHandler 取消关注
func UnfollowHandler(ctx context.Context, c *app.RequestContext) {
	Unfollow(ctx, c)
}

// GetFollowListHandler 获取关注列表
func GetFollowListHandler(ctx context.Context, c *app.RequestContext) {
	GetFollowList(ctx, c)
}

// GetFollowerListHandler 获取粉丝列表
func GetFollowerListHandler(ctx context.Context, c *app.RequestContext) {
	GetFollowerList(ctx, c)
}

// CheckFollowStatusHandler 检查关注状态
func CheckFollowStatusHandler(ctx context.Context, c *app.RequestContext) {
	CheckFollowStatus(ctx, c)
}

// IsFollowingHandler 检查是否关注
func IsFollowingHandler(ctx context.Context, c *app.RequestContext) {
	IsFollowing(ctx, c)
}

// GetFollowCountHandler 获取关注数量
func GetFollowCountHandler(ctx context.Context, c *app.RequestContext) {
	GetFollowCount(ctx, c)
}

// GetFollowerCountHandler 获取粉丝数量
func GetFollowerCountHandler(ctx context.Context, c *app.RequestContext) {
	GetFollowerCount(ctx, c)
}

// GetMutualFollowsHandler 获取共同关注
func GetMutualFollowsHandler(ctx context.Context, c *app.RequestContext) {
	GetMutualFollows(ctx, c)
}