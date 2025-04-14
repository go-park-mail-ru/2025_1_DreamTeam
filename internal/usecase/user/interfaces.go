package usecase

import (
	"context"
	"mime/multipart"
	usermodels "skillForce/internal/models/user"
)

type UserAuthRepository interface {
	RegisterUser(ctx context.Context, user *usermodels.User) (string, error)
	AuthenticateUser(ctx context.Context, email, password string) (string, error)
	GetUserByToken(ctx context.Context, token string) (*usermodels.User, error)
	GetUserByCookie(ctx context.Context, cookieValue string) (*usermodels.UserProfile, error)
	ValidUser(ctx context.Context, user *usermodels.User) (string, error)
	LogoutUser(ctx context.Context, userId int) error
}

type UserProfileRepository interface {
	UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	UpdateProfile(ctx context.Context, userId int, userProfile *usermodels.UserProfile) error
	UpdateProfilePhoto(ctx context.Context, url string, userId int) (string, error)
	DeleteProfilePhoto(ctx context.Context, userId int) error
}

type UserMailRepository interface {
	SendRegMail(ctx context.Context, user *usermodels.User, token string) error
	SendWelcomeMail(ctx context.Context, user *usermodels.User) error
}
