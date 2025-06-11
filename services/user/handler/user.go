package handler

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"hupu/kitex_gen/user"
	"hupu/services/user/repository"
	"hupu/shared/constants"
	"hupu/shared/log"
	"hupu/shared/models"
	"hupu/shared/utils"
)

type UserHandler struct {
	db  repository.UserRepository
	rdb repository.UserRepository
}

func NewUserHandler(db *gorm.DB, rdb *redis.Client) *UserHandler {
	return &UserHandler{
		db:  repository.NewUserRepository(db),
		rdb: repository.NewUserRedisRepo(rdb),
	}
}

func (h *UserHandler) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	logger := log.GetLogger().WithField(constants.TraceIdKey, ctx.Value(constants.TraceIdKey).(string))
	logger.Info("Register start")
	userModel := &models.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Username,
		Phone:    req.Phone,
	}
	savedUser, err := h.rdb.CreateUser(ctx, userModel)
	if err != nil {
		return &user.RegisterResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to create user, err:%s", err.Error()),
		}, nil
	}
	logger.Infof("Register success, user:%+v", savedUser)
	return &user.RegisterResponse{
		Code:    200,
		Message: "注册成功",
		User: &user.User{
			Id:        savedUser.ID,
			Username:  savedUser.Username,
			Nickname:  savedUser.Nickname,
			Avatar:    savedUser.Avatar,
			Phone:     savedUser.Phone,
			Email:     savedUser.Email,
			CreatedAt: savedUser.CreatedAt.Unix(),
			UpdatedAt: savedUser.UpdatedAt.Unix(),
		},
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	getUser, err := h.rdb.GetUser(ctx, &models.User{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return &user.LoginResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to login, err:%s", err.Error()),
		}, nil
	}

	// 生成JWT token
	token, err := utils.GenerateToken(getUser.ID, getUser.Username)
	if err != nil {
		return &user.LoginResponse{
			Code:    500,
			Message: "生成token失败",
		}, err
	}

	return &user.LoginResponse{
		Code:    200,
		Message: "登录成功",
		Token:   token,
		User: &user.User{
			Id:        getUser.ID,
			Username:  getUser.Username,
			Nickname:  getUser.Nickname,
			Avatar:    getUser.Avatar,
			Phone:     getUser.Phone,
			Email:     getUser.Email,
			CreatedAt: getUser.CreatedAt.Unix(),
			UpdatedAt: getUser.UpdatedAt.Unix(),
		},
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	getUser, err := h.rdb.GetUser(ctx, &models.User{
		ID: req.UserId,
	})
	if err != nil {
		return &user.GetUserResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to get user, err:%s", err.Error()),
		}, err
	}

	return &user.GetUserResponse{
		Code:    0,
		Message: "查询成功",
		User: &user.User{
			Id:        getUser.ID,
			Username:  getUser.Username,
			Nickname:  getUser.Nickname,
			Avatar:    getUser.Avatar,
			Phone:     getUser.Phone,
			Email:     getUser.Email,
			CreatedAt: getUser.CreatedAt.Unix(),
			UpdatedAt: getUser.UpdatedAt.Unix(),
		},
	}, nil
}
