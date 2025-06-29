package common

import "github.com/cloudwego/hertz/pkg/app"

type ServiceCall = func() (any, error)

// CallService 统一调用服务的函数
func CallService(c *app.RequestContext, serviceCall ServiceCall, callName string, errorMsg string) {
	resp, err := serviceCall()
	if err != nil {
		HandleServiceError(c, callName, err, errorMsg)
		return
	}
	RespondWithSuccess(c, resp)
}
