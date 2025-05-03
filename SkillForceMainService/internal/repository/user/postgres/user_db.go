package postgres

import (
	"context"
	"fmt"
	"html"
	usermodels "skillForce/internal/models/user"
	"skillForce/pkg/logs"
)

func (d *Database) GetUserByCookie(ctx context.Context, cookieValue string) (*usermodels.UserProfile, error) {
	var userProfile usermodels.UserProfile
	var role string
	err := d.conn.QueryRow("SELECT u.id, u.email, u.name, COALESCE(u.bio, ''), u.avatar_src, u.hide_email, u.role FROM usertable u JOIN sessions s ON u.id = s.user_id WHERE s.token = $1 AND s.expire > NOW();",
		cookieValue).Scan(&userProfile.Id, &userProfile.Email, &userProfile.Name, &userProfile.Bio, &userProfile.AvatarSrc, &userProfile.HideEmail, &role)
	if err != nil {
		logs.PrintLog(ctx, "GetUserByCookie", fmt.Sprintf("error in GetUserByCookie %+v", err))
		return nil, err
	}
	logs.PrintLog(ctx, "GetUserByCookie", fmt.Sprintf("role: %+v", role))
	if role == "admin" {
		userProfile.IsAdmin = true
	}
	userProfile.Email = html.EscapeString(userProfile.Email)
	userProfile.Name = html.EscapeString(userProfile.Name)
	userProfile.Bio = html.EscapeString(userProfile.Bio)
	userProfile.AvatarSrc = html.EscapeString(userProfile.AvatarSrc)
	return &userProfile, err
}

func (d *Database) LogoutUser(ctx context.Context, userId int) error {
	_, err := d.conn.Exec("DELETE FROM sessions WHERE user_id = $1", userId)
	if err != nil {
		return err
	}
	logs.PrintLog(ctx, "LogoutUser", fmt.Sprintf("logout user with id %+v in db", userId))
	return err
}
