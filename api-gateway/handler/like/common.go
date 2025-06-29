package like

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"hupu/api-gateway/handler"
	"hupu/api-gateway/handler/common"
	"hupu/kitex_gen/like/likeservice"
	"hupu/shared/log"
)

// 所有公共方法都已迁移到 common 包中
// 这里保留一些 like 特有的业务逻辑函数

// 以下函数直接使用 common 包中的函数，保持向后兼容

// 这个文件保留为空，所有公共函数已移至 common 包

// ValidateTargetParams 验证目标参数
func ValidateTargetParams(targetID, targetType string) error {
	if targetID == "" {
		return &ValidationError{Message: MsgTargetIDRequired}
	}
	if targetType == "" {
		return &ValidationError{Message: MsgTargetTypeRequired}
	}
	if targetType != TargetTypePost && targetType != TargetTypeComment {
		return &ValidationError{Message: "目标类型必须是post或comment"}
	}
	return nil
}

// ValidationError 验证错误
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// CallLikeService 统一调用点赞服务的函数
func CallLikeService(ctx context.Context, c *app.RequestContext, serviceCall func() (interface{}, error), errorMsg string) {
	resp, err := serviceCall()
	if err != nil {
		log.GetLogger().Errorf("%s error: %v", errorMsg, err)
		common.RespondWithError(c, HTTPStatusInternalServerError, CodeError, errorMsg)
		return
	}
	common.RespondWithSuccess(c, resp)
}

// GetLikeClient 获取点赞服务客户端
func GetLikeClient() likeservice.Client {
	return handler.GetLikeClient()
}
