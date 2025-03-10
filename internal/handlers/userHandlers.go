package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"skillForce/internal/models"
	"skillForce/internal/response"
	"skillForce/internal/usecase"
	"strconv"
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
func setCookie(w http.ResponseWriter, userId int) {
	stingUserId := strconv.Itoa(userId)
	expiration := time.Now().Add(10 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    stingUserId,
		Expires:  expiration,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// SameSite: http.SameSiteNoneMode,
		// Secure:   false,
	}
	http.SetCookie(w, &cookie)
}

// deleteCookie - удаление куки у пользователь
func deleteCookie(w http.ResponseWriter, userId int) {
	stingUserId := strconv.Itoa(userId)
	expiration := time.Now().AddDate(0, 0, -1)
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    stingUserId,
		Expires:  expiration,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
}

// checkCookie - проверка наличия куки
func (h *UserHandler) checkCookie(r *http.Request) *models.User {
	session, err := r.Cookie("session_id")
	loggedIn := (err != http.ErrNoCookie)
	if loggedIn {
		user, err := h.useCase.GetUserByCookie(session.Value)
		if err != nil {
			log.Print(err)
			return nil
		}
		return user
	}
	return nil
}

func (h *UserHandler) IsAuthorized(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if session == nil {
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}
	loggedIn := (err != http.ErrNoCookie)
	if loggedIn {
		user, err := h.useCase.GetUserByCookie(session.Value)
		if err != nil {
			log.Print(err)
			response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
			return
		}
		if user != nil {
			response.SendUser(user, w, r)
			return
		}
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}
	response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
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
	userFromCoockies := h.checkCookie(r)
	if userFromCoockies != nil {
		log.Print("user already registered in")
		response.SendOKResponse(w, r)
		return
	}

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

	err = h.useCase.RegisterUser(&user)
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

	setCookie(w, user.Id)
	response.SendOKResponse(w, r)
}

// LoginUser - обработчик авторизации пользователя
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	userFromCoockies := h.checkCookie(r)
	if userFromCoockies != nil {
		log.Print("user already logged in")
		response.SendOKResponse(w, r)
		return
	}

	if r.Method == http.MethodOptions {
		response.SendCors(w, r)
		return
	}

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

	userId, err := h.useCase.AuthenticateUser(&user)
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

	setCookie(w, userId)
	response.SendOKResponse(w, r)
}

// LogoutUser - обработчик для выхода из сессии у пользователя
func (h *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	userFromCoockies := h.checkCookie(r)
	if userFromCoockies != nil {
		log.Printf("logout user %+v", userFromCoockies)
		err := h.useCase.LogoutUser(userFromCoockies.Id)
		if err != nil {
			log.Printf("from logoutUser: %v", err) //TODO: тут дырка, надо подумать, че делать...
		}
		deleteCookie(w, userFromCoockies.Id)

	}
	response.SendOKResponse(w, r)
}
