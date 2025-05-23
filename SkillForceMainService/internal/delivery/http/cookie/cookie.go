package cookie

import (
	"context"
	"fmt"
	"net/http"
	models "skillForce/internal/models/user"
	"skillForce/pkg/logs"
	"time"
)

type CookieUsecaseInterface interface {
	GetUserByCookie(ctx context.Context, cookieValue string) (*models.UserProfile, error)
	LogoutUser(ctx context.Context, userId int) error
}

type CookieManager struct {
	usecase CookieUsecaseInterface
}

func NewCookieManager(usecase CookieUsecaseInterface) *CookieManager {
	return &CookieManager{usecase: usecase}
}

func (c *CookieManager) SetCookie(w http.ResponseWriter, cookieValue string) {
	expiration := time.Now().AddDate(1, 0, 0)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    cookieValue,
		Expires:  expiration,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// SameSite: http.SameSiteNoneMode,
		// Secure:   false,
	}
	http.SetCookie(w, &cookie)
}

func (c *CookieManager) DeleteCookie(w http.ResponseWriter) {
	cookieValue := "hello from server"
	expiration := time.Now().AddDate(0, 0, -1)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    cookieValue,
		Expires:  expiration,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
}

func (c *CookieManager) CheckCookie(r *http.Request) *models.UserProfile {
	session, err := r.Cookie("session_id")
	logs.PrintLog(r.Context(), "checkCookie", "checking cookie")
	loggedIn := (err != http.ErrNoCookie)
	if loggedIn {
		userProfile, err := c.usecase.GetUserByCookie(r.Context(), session.Value)
		if err != nil {
			logs.PrintLog(r.Context(), "checkCookie", fmt.Sprintf("%+v", err))
			return nil
		}
		return userProfile
	}
	return nil
}

func (c *CookieManager) LogoutUser(ctx context.Context, userId int) error {
	return c.usecase.LogoutUser(ctx, userId)
}
