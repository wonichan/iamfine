package handler

import (
	"github.com/cloudwego/kitex/client"

	"hupu/kitex_gen/comment/commentservice"
	"hupu/kitex_gen/follow/followservice"
	"hupu/kitex_gen/like/likeservice"
	"hupu/kitex_gen/notification/notificationservice"
	"hupu/kitex_gen/post/postservice"
	"hupu/kitex_gen/user/userservice"
	"hupu/shared/config"
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

func Init() {
	userClient = userservice.MustNewClient(userServiceName, client.WithHostPorts(config.GlobalConfig.Services.User.Host+":"+config.GlobalConfig.Services.User.Port))
	postClient = postservice.MustNewClient(postServiceName, client.WithHostPorts(config.GlobalConfig.Services.Post.Host+":"+config.GlobalConfig.Services.Post.Port))
	notificationClient = notificationservice.MustNewClient(notificationName, client.WithHostPorts(config.GlobalConfig.Services.Notification.Host+":"+config.GlobalConfig.Services.Notification.Port))
	likeClient = likeservice.MustNewClient(likeServiceName, client.WithHostPorts(config.GlobalConfig.Services.Like.Host+":"+config.GlobalConfig.Services.Like.Port))
	followClient = followservice.MustNewClient(followServiceName, client.WithHostPorts(config.GlobalConfig.Services.Follow.Host+":"+config.GlobalConfig.Services.Follow.Port))
	commentClient = commentservice.MustNewClient(commentServiceName, client.WithHostPorts(config.GlobalConfig.Services.Comment.Host+":"+config.GlobalConfig.Services.Comment.Port))
}
