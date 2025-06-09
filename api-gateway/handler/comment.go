package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"hupu/kitex_gen/comment"
	"hupu/kitex_gen/comment/commentservice"
	"hupu/shared/config"
)

// 创建评论
func CreateComment(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	// 解析请求参数
	var req comment.CreateCommentRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 设置用户ID
	req.UserId = userID.(string)

	// 创建评论服务客户端
	client, err := commentservice.NewClient("comment", client.WithHostPorts(config.GlobalConfig.Services.Comment.Host+":"+config.GlobalConfig.Services.Comment.Port))
	if err != nil {
		hlog.Errorf("Create comment client error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "服务连接失败",
		})
		return
	}

	// 调用评论服务
	resp, err := client.CreateComment(ctx, &req)
	if err != nil {
		hlog.Errorf("CreateComment error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建评论失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// 获取评论列表
func GetCommentList(ctx context.Context, c *app.RequestContext) {
	// 获取帖子ID参数
	postID := c.Param("post_id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "帖子ID不能为空",
		})
		return
	}

	// 解析查询参数
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 创建评论服务客户端
	client, err := commentservice.NewClient("comment", client.WithHostPorts(config.GlobalConfig.Services.Comment.Host+":"+config.GlobalConfig.Services.Comment.Port))
	if err != nil {
		hlog.Errorf("Create comment client error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "服务连接失败",
		})
		return
	}

	// 调用评论服务
	req := &comment.GetCommentListRequest{
		PostId:   postID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}
	resp, err := client.GetCommentList(ctx, req)
	if err != nil {
		hlog.Errorf("GetCommentList error: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取评论列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
