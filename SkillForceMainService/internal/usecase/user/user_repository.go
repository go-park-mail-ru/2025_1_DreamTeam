package usecase

import (
	context "context"
	usermodels "skillForce/internal/models/user"
)

type UserRepository interface {
	GetUserByCookie(ctx context.Context, cookieValue string) (*usermodels.UserProfile, error)
	LogoutUser(ctx context.Context, userId int) error
}
