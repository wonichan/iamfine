package common

import "github.com/cloudwego/hertz/pkg/app"

type RpcCall interface {
	Call() (any, error)
}

type ServiceCall func() (any, error)

func (s ServiceCall) Call() (any, error) {
	return s()
}

// CallService 统一调用服务的函数
func CallService(c *app.RequestContext, rc RpcCall, callName string, errorMsg string) {
	resp, err := rc.Call()
	if err != nil {
		HandleServiceError(c, callName, err, errorMsg)
		return
	}
	RespondWithSuccess(c, resp)
}
