package services

import (
	"context"

	"github.com/besanh/mini-crm/models"
	"github.com/besanh/mini-crm/repositories"
	"github.com/google/uuid"
)

type (
	IUsers interface {
		GetUserByID(ctx context.Context, id uuid.UUID) (result *models.UserResponse, err error)
	}

	Users struct {
		userRepo repositories.IUsers
	}
)

func NewUsers(userRepo repositories.IUsers) IUsers {
	return &Users{
		userRepo: userRepo,
	}
}

func (u *Users) GetUserByID(ctx context.Context, id uuid.UUID) (result *models.UserResponse, err error) {
	return u.userRepo.GetUserByID(ctx, id)
}
