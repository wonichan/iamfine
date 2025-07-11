// Code generated by Kitex v0.12.3. DO NOT EDIT.

package commentservice

import (
	"context"
	"errors"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	comment "hupu/kitex_gen/comment"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"CreateComment": kitex.NewMethodInfo(
		createCommentHandler,
		newCommentServiceCreateCommentArgs,
		newCommentServiceCreateCommentResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetCommentList": kitex.NewMethodInfo(
		getCommentListHandler,
		newCommentServiceGetCommentListArgs,
		newCommentServiceGetCommentListResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetComment": kitex.NewMethodInfo(
		getCommentHandler,
		newCommentServiceGetCommentArgs,
		newCommentServiceGetCommentResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetUserComments": kitex.NewMethodInfo(
		getUserCommentsHandler,
		newCommentServiceGetUserCommentsArgs,
		newCommentServiceGetUserCommentsResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"DeleteComment": kitex.NewMethodInfo(
		deleteCommentHandler,
		newCommentServiceDeleteCommentArgs,
		newCommentServiceDeleteCommentResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"RateComment": kitex.NewMethodInfo(
		rateCommentHandler,
		newCommentServiceRateCommentArgs,
		newCommentServiceRateCommentResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetUserCommentRating": kitex.NewMethodInfo(
		getUserCommentRatingHandler,
		newCommentServiceGetUserCommentRatingArgs,
		newCommentServiceGetUserCommentRatingResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"UpdateCommentRating": kitex.NewMethodInfo(
		updateCommentRatingHandler,
		newCommentServiceUpdateCommentRatingArgs,
		newCommentServiceUpdateCommentRatingResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"DeleteCommentRating": kitex.NewMethodInfo(
		deleteCommentRatingHandler,
		newCommentServiceDeleteCommentRatingArgs,
		newCommentServiceDeleteCommentRatingResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
}

var (
	commentServiceServiceInfo                = NewServiceInfo()
	commentServiceServiceInfoForClient       = NewServiceInfoForClient()
	commentServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return commentServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return commentServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return commentServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "CommentService"
	handlerType := (*comment.CommentService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "comment",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.12.3",
		Extra:           extra,
	}
	return svcInfo
}

func createCommentHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*comment.CommentServiceCreateCommentArgs)
	realResult := result.(*comment.CommentServiceCreateCommentResult)
	success, err := handler.(comment.CommentService).CreateComment(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCommentServiceCreateCommentArgs() interface{} {
	return comment.NewCommentServiceCreateCommentArgs()
}

func newCommentServiceCreateCommentResult() interface{} {
	return comment.NewCommentServiceCreateCommentResult()
}

func getCommentListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*comment.CommentServiceGetCommentListArgs)
	realResult := result.(*comment.CommentServiceGetCommentListResult)
	success, err := handler.(comment.CommentService).GetCommentList(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCommentServiceGetCommentListArgs() interface{} {
	return comment.NewCommentServiceGetCommentListArgs()
}

func newCommentServiceGetCommentListResult() interface{} {
	return comment.NewCommentServiceGetCommentListResult()
}

func getCommentHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*comment.CommentServiceGetCommentArgs)
	realResult := result.(*comment.CommentServiceGetCommentResult)
	success, err := handler.(comment.CommentService).GetComment(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCommentServiceGetCommentArgs() interface{} {
	return comment.NewCommentServiceGetCommentArgs()
}

func newCommentServiceGetCommentResult() interface{} {
	return comment.NewCommentServiceGetCommentResult()
}

func getUserCommentsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*comment.CommentServiceGetUserCommentsArgs)
	realResult := result.(*comment.CommentServiceGetUserCommentsResult)
	success, err := handler.(comment.CommentService).GetUserComments(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCommentServiceGetUserCommentsArgs() interface{} {
	return comment.NewCommentServiceGetUserCommentsArgs()
}

func newCommentServiceGetUserCommentsResult() interface{} {
	return comment.NewCommentServiceGetUserCommentsResult()
}

func deleteCommentHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*comment.CommentServiceDeleteCommentArgs)
	realResult := result.(*comment.CommentServiceDeleteCommentResult)
	success, err := handler.(comment.CommentService).DeleteComment(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCommentServiceDeleteCommentArgs() interface{} {
	return comment.NewCommentServiceDeleteCommentArgs()
}

func newCommentServiceDeleteCommentResult() interface{} {
	return comment.NewCommentServiceDeleteCommentResult()
}

func rateCommentHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*comment.CommentServiceRateCommentArgs)
	realResult := result.(*comment.CommentServiceRateCommentResult)
	success, err := handler.(comment.CommentService).RateComment(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCommentServiceRateCommentArgs() interface{} {
	return comment.NewCommentServiceRateCommentArgs()
}

func newCommentServiceRateCommentResult() interface{} {
	return comment.NewCommentServiceRateCommentResult()
}

func getUserCommentRatingHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*comment.CommentServiceGetUserCommentRatingArgs)
	realResult := result.(*comment.CommentServiceGetUserCommentRatingResult)
	success, err := handler.(comment.CommentService).GetUserCommentRating(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCommentServiceGetUserCommentRatingArgs() interface{} {
	return comment.NewCommentServiceGetUserCommentRatingArgs()
}

func newCommentServiceGetUserCommentRatingResult() interface{} {
	return comment.NewCommentServiceGetUserCommentRatingResult()
}

func updateCommentRatingHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*comment.CommentServiceUpdateCommentRatingArgs)
	realResult := result.(*comment.CommentServiceUpdateCommentRatingResult)
	success, err := handler.(comment.CommentService).UpdateCommentRating(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCommentServiceUpdateCommentRatingArgs() interface{} {
	return comment.NewCommentServiceUpdateCommentRatingArgs()
}

func newCommentServiceUpdateCommentRatingResult() interface{} {
	return comment.NewCommentServiceUpdateCommentRatingResult()
}

func deleteCommentRatingHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*comment.CommentServiceDeleteCommentRatingArgs)
	realResult := result.(*comment.CommentServiceDeleteCommentRatingResult)
	success, err := handler.(comment.CommentService).DeleteCommentRating(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCommentServiceDeleteCommentRatingArgs() interface{} {
	return comment.NewCommentServiceDeleteCommentRatingArgs()
}

func newCommentServiceDeleteCommentRatingResult() interface{} {
	return comment.NewCommentServiceDeleteCommentRatingResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) CreateComment(ctx context.Context, req *comment.CreateCommentRequest) (r *comment.CreateCommentResponse, err error) {
	var _args comment.CommentServiceCreateCommentArgs
	_args.Req = req
	var _result comment.CommentServiceCreateCommentResult
	if err = p.c.Call(ctx, "CreateComment", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetCommentList(ctx context.Context, req *comment.GetCommentListRequest) (r *comment.GetCommentListResponse, err error) {
	var _args comment.CommentServiceGetCommentListArgs
	_args.Req = req
	var _result comment.CommentServiceGetCommentListResult
	if err = p.c.Call(ctx, "GetCommentList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetComment(ctx context.Context, req *comment.GetCommentRequest) (r *comment.GetCommentResponse, err error) {
	var _args comment.CommentServiceGetCommentArgs
	_args.Req = req
	var _result comment.CommentServiceGetCommentResult
	if err = p.c.Call(ctx, "GetComment", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetUserComments(ctx context.Context, req *comment.GetUserCommentsRequest) (r *comment.GetUserCommentsResponse, err error) {
	var _args comment.CommentServiceGetUserCommentsArgs
	_args.Req = req
	var _result comment.CommentServiceGetUserCommentsResult
	if err = p.c.Call(ctx, "GetUserComments", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DeleteComment(ctx context.Context, req *comment.DeleteCommentRequest) (r *comment.DeleteCommentResponse, err error) {
	var _args comment.CommentServiceDeleteCommentArgs
	_args.Req = req
	var _result comment.CommentServiceDeleteCommentResult
	if err = p.c.Call(ctx, "DeleteComment", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) RateComment(ctx context.Context, req *comment.RateCommentRequest) (r *comment.RateCommentResponse, err error) {
	var _args comment.CommentServiceRateCommentArgs
	_args.Req = req
	var _result comment.CommentServiceRateCommentResult
	if err = p.c.Call(ctx, "RateComment", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetUserCommentRating(ctx context.Context, req *comment.GetUserCommentRatingRequest) (r *comment.GetUserCommentRatingResponse, err error) {
	var _args comment.CommentServiceGetUserCommentRatingArgs
	_args.Req = req
	var _result comment.CommentServiceGetUserCommentRatingResult
	if err = p.c.Call(ctx, "GetUserCommentRating", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) UpdateCommentRating(ctx context.Context, req *comment.UpdateCommentRatingRequest) (r *comment.UpdateCommentRatingResponse, err error) {
	var _args comment.CommentServiceUpdateCommentRatingArgs
	_args.Req = req
	var _result comment.CommentServiceUpdateCommentRatingResult
	if err = p.c.Call(ctx, "UpdateCommentRating", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) DeleteCommentRating(ctx context.Context, req *comment.DeleteCommentRatingRequest) (r *comment.DeleteCommentRatingResponse, err error) {
	var _args comment.CommentServiceDeleteCommentRatingArgs
	_args.Req = req
	var _result comment.CommentServiceDeleteCommentRatingResult
	if err = p.c.Call(ctx, "DeleteCommentRating", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
