package repository

import (
	"context"

	"hupu/shared/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, models *models.User) (*models.User, error)
	GetUser(ctx context.Context, models *models.User) (*models.User, error)
	GetUserByUsername(ctx context.Context, models *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, models *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, models *models.User) (*models.User, error)
}
