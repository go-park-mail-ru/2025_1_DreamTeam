package userinfrastructure

import (
	"context"
	"mime/multipart"
	usermodels "skillForce/internal/models/user"
	"skillForce/internal/repository/mail"
	"skillForce/internal/repository/minio"
	"skillForce/internal/repository/postgres"
)

type UserInfrastructure struct {
	Database *postgres.Database
	Minio    *minio.Minio
	Mail     *mail.Mail
}

func NewUserInfrastructure(db *postgres.Database, mail *mail.Mail, minio *minio.Minio) *UserInfrastructure {
	return &UserInfrastructure{
		Database: db,
		Minio:    minio,
		Mail:     mail,
	}
}

func (u *UserInfrastructure) RegisterUser(ctx context.Context, user *usermodels.User) (string, error) {
	return u.Database.RegisterUser(ctx, user)
}

func (u *UserInfrastructure) AuthenticateUser(ctx context.Context, email, password string) (string, error) {
	return u.Database.AuthenticateUser(ctx, email, password)
}

func (i *UserInfrastructure) GetUserByCookie(ctx context.Context, cookieValue string) (*usermodels.UserProfile, error) {
	return i.Database.GetUserByCookie(ctx, cookieValue)
}

func (i *UserInfrastructure) LogoutUser(ctx context.Context, userId int) error {
	return i.Database.LogoutUser(ctx, userId)
}

func (i *UserInfrastructure) UpdateProfile(ctx context.Context, userId int, userProfile *usermodels.UserProfile) error {
	return i.Database.UpdateProfile(ctx, userId, userProfile)
}

func (i *UserInfrastructure) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	return i.Minio.UploadFileToMinIO(ctx, file, fileHeader)
}

func (i *UserInfrastructure) UpdateProfilePhoto(ctx context.Context, photo_url string, userId int) (string, error) {
	return i.Database.UpdateProfilePhoto(ctx, photo_url, userId)
}

func (i *UserInfrastructure) DeleteProfilePhoto(ctx context.Context, userId int) error {
	return i.Database.DeleteProfilePhoto(ctx, userId)
}

func (i *UserInfrastructure) ValidUser(ctx context.Context, user *usermodels.User) (string, error) {
	return i.Database.ValidUser(ctx, user)
}

func (i *UserInfrastructure) GetUserByToken(ctx context.Context, token string) (*usermodels.User, error) {
	return i.Database.GetUserByToken(ctx, token)
}

func (i *UserInfrastructure) SendRegMail(ctx context.Context, user *usermodels.User, token string) error {
	return i.Mail.SendRegMail(ctx, user, token)
}

func (i *UserInfrastructure) SendWelcomeMail(ctx context.Context, user *usermodels.User) error {
	return i.Mail.SendWelcomeMail(ctx, user)
}
