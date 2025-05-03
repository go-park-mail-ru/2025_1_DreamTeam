package userinfrastructure

import (
	"context"
	"fmt"
	"log"
	"skillForce/config"
	usermodels "skillForce/internal/models/user"
	"skillForce/internal/repository/user/postgres"
)

type UserInfrastructure struct {
	Database *postgres.Database
}

func NewUserInfrastructure(conf *config.Config) *UserInfrastructure {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Name)
	database, err := postgres.NewDatabase(dsn, conf.Secrets.JwtSessionSecret)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return &UserInfrastructure{
		Database: database,
	}
}

func (u *UserInfrastructure) Close() {
	u.Database.Close()
}

func (u *UserInfrastructure) GetUserByCookie(ctx context.Context, cookieValue string) (*usermodels.UserProfile, error) {
	return u.Database.GetUserByCookie(ctx, cookieValue)
}

func (u *UserInfrastructure) LogoutUser(ctx context.Context, userId int) error {
	return u.Database.LogoutUser(ctx, userId)
}
