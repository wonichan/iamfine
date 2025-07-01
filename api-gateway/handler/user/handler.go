package user

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// Register 用户注册
func RegisterHandler(ctx context.Context, c *app.RequestContext) {
	Register(ctx, c)
}

// Login 用户登录
func LoginHandler(ctx context.Context, c *app.RequestContext) {
	Login(ctx, c)
}

// GetUser 获取用户信息
func GetUserHandler(ctx context.Context, c *app.RequestContext) {
	GetUser(ctx, c)
}

// UpdateUser 更新用户信息
func UpdateUserHandler(ctx context.Context, c *app.RequestContext) {
	UpdateUser(ctx, c)
}

// FollowUser 关注用户
func FollowUserHandler(ctx context.Context, c *app.RequestContext) {
	FollowUser(ctx, c)
}

// UnfollowUser 取消关注用户
func UnfollowUserHandler(ctx context.Context, c *app.RequestContext) {
	UnfollowUser(ctx, c)
}

// GetFollowers 获取粉丝列表
func GetFollowersHandler(ctx context.Context, c *app.RequestContext) {
	GetFollowers(ctx, c)
}

// GetFollowing 获取关注列表
func GetFollowingHandler(ctx context.Context, c *app.RequestContext) {
	GetFollowing(ctx, c)
}

// CreateAnonymousProfile 创建匿名马甲
func CreateAnonymousProfileHandler(ctx context.Context, c *app.RequestContext) {
	CreateAnonymousProfile(ctx, c)
}

// GetAnonymousProfiles 获取匿名马甲列表
func GetAnonymousProfilesHandler(ctx context.Context, c *app.RequestContext) {
	GetAnonymousProfiles(ctx, c)
}

// UpdateAnonymousProfile 更新匿名马甲
func UpdateAnonymousProfileHandler(ctx context.Context, c *app.RequestContext) {
	UpdateAnonymousProfile(ctx, c)
}

// GetUserStats 获取用户统计信息
func GetUserStatsHandler(ctx context.Context, c *app.RequestContext) {
	GetUserStats(ctx, c)
}

// WxLogin 微信登录
func WxLoginHandler(ctx context.Context, c *app.RequestContext) {
	WxLogin(ctx, c)
}

// GetUserInfo 获取用户信息
func GetUserInfoHandler(ctx context.Context, c *app.RequestContext) {
	GetUserInfo(ctx, c)
}

// GetUnreadCount 获取未读消息数
func GetUnreadCountHandler(ctx context.Context, c *app.RequestContext) {
	GetUnreadCount(ctx, c)
}
