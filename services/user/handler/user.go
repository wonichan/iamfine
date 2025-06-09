package handler

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"gorm.io/gorm"
	"hupu/kitex_gen/user"
	"hupu/shared/models"
	"hupu/shared/utils"
)

type UserHandler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewUserHandler(db *gorm.DB, rdb *redis.Client) *UserHandler {
	return &UserHandler{
		db:  db,
		rdb: rdb,
	}
}

func (h *UserHandler) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	// 检查用户名是否已存在
	var existUser models.User
	err := h.db.Where("username = ? OR phone = ?", req.Username, req.Phone).First(&existUser).Error
	if err == nil {
		return &user.RegisterResponse{
			Code:    400,
			Message: "用户名或手机号已存在",
		}, nil
	}

	// 生成用户ID
	userID := xid.New().String()

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return &user.RegisterResponse{
			Code:    500,
			Message: "密码加密失败",
		}, err
	}

	// 创建用户
	newUser := models.User{
		ID:       userID,
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Username,
		Phone:    req.Phone,
	}

	err = h.db.Create(&newUser).Error
	if err != nil {
		return &user.RegisterResponse{
			Code:    500,
			Message: "创建用户失败",
		}, err
	}

	return &user.RegisterResponse{
		Code:    200,
		Message: "注册成功",
		User: &user.User{
			Id:        newUser.ID,
			Username:  newUser.Username,
			Nickname:  newUser.Nickname,
			Avatar:    newUser.Avatar,
			Phone:     newUser.Phone,
			Email:     newUser.Email,
			CreatedAt: newUser.CreatedAt.Unix(),
			UpdatedAt: newUser.UpdatedAt.Unix(),
		},
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	// 查找用户
	var userModel models.User
	err := h.db.Where("username = ?", req.Username).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.LoginResponse{
				Code:    400,
				Message: "用户不存在",
			}, nil
		}
		return &user.LoginResponse{
			Code:    500,
			Message: "查询用户失败",
		}, err
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, userModel.Password) {
		return &user.LoginResponse{
			Code:    400,
			Message: "密码错误",
		}, nil
	}

	// 生成JWT token
	token, err := utils.GenerateToken(userModel.ID, userModel.Username)
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
			Id:        userModel.ID,
			Username:  userModel.Username,
			Nickname:  userModel.Nickname,
			Avatar:    userModel.Avatar,
			Phone:     userModel.Phone,
			Email:     userModel.Email,
			CreatedAt: userModel.CreatedAt.Unix(),
			UpdatedAt: userModel.UpdatedAt.Unix(),
		},
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	var userModel models.User
	err := h.db.Where("id = ?", req.UserId).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.GetUserResponse{
				Code:    404,
				Message: "用户不存在",
			}, nil
		}
		return &user.GetUserResponse{
			Code:    500,
			Message: "查询用户失败",
		}, err
	}

	return &user.GetUserResponse{
		Code:    200,
		Message: "查询成功",
		User: &user.User{
			Id:        userModel.ID,
			Username:  userModel.Username,
			Nickname:  userModel.Nickname,
			Avatar:    userModel.Avatar,
			Phone:     userModel.Phone,
			Email:     userModel.Email,
			CreatedAt: userModel.CreatedAt.Unix(),
			UpdatedAt: userModel.UpdatedAt.Unix(),
		},
	}, nil
}
