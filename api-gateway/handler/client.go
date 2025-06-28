package handler

import (
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"

	"hupu/kitex_gen/comment/commentservice"
	"hupu/kitex_gen/follow/followservice"
	"hupu/kitex_gen/like/likeservice"
	"hupu/kitex_gen/notification/notificationservice"
	"hupu/kitex_gen/post/postservice"
	"hupu/kitex_gen/user/userservice"
	"hupu/shared/config"
	"hupu/shared/utils"
)

var (
	userClient         userservice.Client
	postClient         postservice.Client
	notificationClient notificationservice.Client
	likeClient         likeservice.Client
	followClient       followservice.Client
	commentClient      commentservice.Client
)

const (
	userServiceName    = "user"
	postServiceName    = "post"
	notificationName   = "notification"
	likeServiceName    = "like"
	followServiceName  = "follow"
	commentServiceName = "comment"
)

func GetUserClient() userservice.Client {
	return userClient
}

func GetPostClient() postservice.Client {
	return postClient
}

func GetNotificationClient() notificationservice.Client {
	return notificationClient
}

func GetLikeClient() likeservice.Client {
	return likeClient
}

func GetFollowClient() followservice.Client {
	return followClient
}

func GetCommentClient() commentservice.Client {
	return commentClient
}

func WithCommonOption(clientName string) []client.Option {
	return []client.Option{
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: clientName}),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		client.WithTracer(utils.NewKitexClientTracer()),
		client.WithTransportProtocol(transport.TTHeader),
		client.WithRPCTimeout(time.Minute),
	}
}

func WithServiceOptions(addr, clientName string) []client.Option {
	return append(
		[]client.Option{client.WithHostPorts(addr)},
		WithCommonOption(clientName)...,
	)
}

func Init() {
	userClient = userservice.MustNewClient(userServiceName,
		WithServiceOptions(config.GlobalConfig.Services.User.Host+":"+config.GlobalConfig.Services.User.Port, "userClient")...)

	postClient = postservice.MustNewClient(postServiceName,
		WithServiceOptions(config.GlobalConfig.Services.Post.Host+":"+config.GlobalConfig.Services.Post.Port, "postClient")...)

	notificationClient = notificationservice.MustNewClient(notificationName,
		WithServiceOptions(config.GlobalConfig.Services.Notification.Host+":"+config.GlobalConfig.Services.Notification.Port, "notificationClient")...)

	likeClient = likeservice.MustNewClient(likeServiceName,
		WithServiceOptions(config.GlobalConfig.Services.Like.Host+":"+config.GlobalConfig.Services.Like.Port, "likeClient")...)

	followClient = followservice.MustNewClient(followServiceName,
		WithServiceOptions(config.GlobalConfig.Services.Follow.Host+":"+config.GlobalConfig.Services.Follow.Port, "followClient")...)

	commentClient = commentservice.MustNewClient(commentServiceName,
		WithServiceOptions(config.GlobalConfig.Services.Comment.Host+":"+config.GlobalConfig.Services.Comment.Port, "commentClient")...)
}
