package handlers

import (
	"fmt"
	"net/http"
	"skillForce/internal/models"
	"skillForce/internal/usecase"
	"skillForce/pkg/logs"
	"time"
)

type Handler struct {
	useCase usecase.UsecaseInterface
}

func NewHandler(uc *usecase.Usecase) *Handler {
	return &Handler{
		useCase: uc,
	}
}

// setCookie - установка куки для контроля, авторизован ли пользователь
func setCookie(w http.ResponseWriter, cookieValue string) {
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

// deleteCookie - удаление куки у пользователь
func deleteCookie(w http.ResponseWriter) {
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

// checkCookie - проверка наличия куки
func (h *Handler) checkCookie(r *http.Request) *models.UserProfile {
	session, err := r.Cookie("session_id")
	logs.PrintLog(r.Context(), "checkCookie", "checking cookie")
	loggedIn := (err != http.ErrNoCookie)
	if loggedIn {
		userProfile, err := h.useCase.GetUserByCookie(r.Context(), session.Value)
		if err != nil {
			logs.PrintLog(r.Context(), "checkCookie", fmt.Sprintf("%+v", err))
			return nil
		}
		return userProfile
	}
	return nil
}
