package repositories

import (
	"context"

	"github.com/besanh/mini-crm/common/log"
	"github.com/besanh/mini-crm/ent"
	"github.com/besanh/mini-crm/ent/users"
	"github.com/besanh/mini-crm/models"
	"github.com/google/uuid"
)

type (
	IUsers interface {
		GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error)
	}
	Users struct {
		client *ent.Client
	}
)

func NewUsers(client *ent.Client) IUsers {
	return &Users{
		client: client,
	}
}

func convertToUserResponse(u *ent.Users) *models.UserResponse {
	return &models.UserResponse{
		GBase: &models.GBase{
			Id:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		UserProfile: models.UserProfile{
			Name:          u.UserProfile.Name,
			Email:         u.UserProfile.Email,
			Sub:           u.UserProfile.Sub,
			Picture:       u.UserProfile.Picture,
			Locale:        u.UserProfile.Locale,
			Profile:       u.UserProfile.Profile,
			GivenName:     u.UserProfile.GivenName,
			FamilyName:    u.UserProfile.FamilyName,
			EmailVerified: u.UserProfile.EmailVerified,
		},
		Status: u.Status.String(),
		Scope:  u.Scope,
	}
}

func (r *Users) GetUserByID(ctx context.Context, id uuid.UUID) (result *models.UserResponse, err error) {
	res, err := r.client.Users.Query().
		Where(
			users.ID(id),
		).
		Only(ctx)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	result = convertToUserResponse(res)
	return
}
