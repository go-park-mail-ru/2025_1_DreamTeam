package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"skillForce/internal/models"
	"skillForce/internal/response"
	"skillForce/internal/usecase"
	"time"

	"github.com/badoux/checkmail"
)

// UserHandler - структура обработчика HTTP-запросов
type UserHandler struct {
	useCase *usecase.UserUsecase
}

// NewUserHandler - конструктор
func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{useCase: uc}
}

// setCookie - установка куки для контроля, авторизован ли пользователь
func setCookie(w http.ResponseWriter, cookieValue string) {
	expiration := time.Now().Add(10 * time.Hour)
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
func (h *UserHandler) checkCookie(r *http.Request) *models.UserProfile {
	session, err := r.Cookie("session_id")
	loggedIn := (err != http.ErrNoCookie)
	if loggedIn {
		userProfile, err := h.useCase.GetUserByCookie(session.Value)
		if err != nil {
			log.Print(err)
			return nil
		}
		return userProfile
	}
	return nil
}

// isValidRegistrationFields - валидация полей регистрации
func isValidRegistrationFields(user *models.User) error {
	if user.Name == "" {
		return errors.New("missing required fields") //TODO: сделать другую ошибку, но пока так, чтобы не поломать фронт
	}

	if len(user.Password) < 5 {
		return errors.New("password too short")
	}

	err := checkmail.ValidateFormat(user.Email) //TODO: улучшить проверку почты
	if err != nil {
		return errors.New("invalid email")
	}

	return nil
}

// isValidLoginFields - валидация полей авторизации
func isValidLoginFields(user *models.User) error {
	if len(user.Password) < 5 {
		return errors.New("password too short")
	}

	err := checkmail.ValidateFormat(user.Email) //TODO: улучшить проверку почты
	if err != nil {
		return errors.New("invalid email")
	}

	return nil
}

// RegisterUser - обработчик регистрации пользователя
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Print("from registerUser: method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("from registerUser: %v", err)
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	err = isValidRegistrationFields(&user)
	if err != nil {
		log.Printf("from registerUser: %v", err)
		response.SendErrorResponse(err.Error(), http.StatusBadRequest, w, r)
		return
	}

	cookie, err := h.useCase.RegisterUser(&user)
	if err != nil {
		log.Printf("from registerUser: %v", err)

		if err.Error() == "email exists" { //TODO: хорошо бы все константы вывести в отдельный файл
			response.SendErrorResponse(err.Error(), http.StatusNotFound, w, r)
			return
		}

		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}

	log.Printf("user %v registered, send him cookie", user)

	setCookie(w, cookie)
	response.SendOKResponse(w, r)
}

// LoginUser - обработчик авторизации пользователя
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Print("from loginUser: method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("from loginUser: %v", err)
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	err = isValidLoginFields(&user)
	if err != nil {
		log.Printf("from loginUser: %v", err)
		response.SendErrorResponse(err.Error(), http.StatusBadRequest, w, r)
		return
	}

	cookieValue, err := h.useCase.AuthenticateUser(&user)
	if err != nil {
		log.Printf("from loginUser: %v", err)

		if err.Error() == "email or password incorrect" { //TODO: хорошо бы все константы вывести в отдельный файл
			response.SendErrorResponse(err.Error(), http.StatusNotFound, w, r)
			return
		}

		response.SendErrorResponse("server error", http.StatusNotFound, w, r)
		return
	}

	log.Printf("user %v login, send him cookie", user)

	setCookie(w, cookieValue)
	response.SendOKResponse(w, r)
}

// LogoutUser - обработчик для выхода из сессии у пользователя
func (h *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	userProfile := h.checkCookie(r)
	if userProfile != nil {
		log.Printf("logout user %+v", userProfile)
		err := h.useCase.LogoutUser(userProfile.Id)
		if err != nil {
			log.Printf("from logoutUser: %v", err) //TODO: тут дырка, надо подумать, че делать...
		}
		deleteCookie(w)

	}
	response.SendOKResponse(w, r)
}

// IsAuthorized - обработчик проверки авторизации, если пользователь авторизован, то возвращает его данные
func (h *UserHandler) IsAuthorized(w http.ResponseWriter, r *http.Request) {
	userProfile := h.checkCookie(r)
	if userProfile != nil {
		response.SendUserProfile(userProfile, w, r)
		return
	}
	log.Print("from isAuthorized: user not logged in")
	response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userProfile := h.checkCookie(r)
	if userProfile == nil {
		log.Print("from updateProfile: user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	var newUserProfile models.UserProfile
	err := json.NewDecoder(r.Body).Decode(&newUserProfile)
	if err != nil {
		log.Printf("from updateProfile: %v", err)
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	//TODO: добавить тут валидацию
	err = h.useCase.UpdateProfile(userProfile.Id, &newUserProfile)
	if err != nil {
		log.Printf("from updateProfile: %v", err)
		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}

	log.Printf("user %v updated profile with values %+v", userProfile, newUserProfile)

	response.SendOKResponse(w, r)
}
