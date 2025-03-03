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
	}
	http.SetCookie(w, &cookie)
}

// checkCookie - проверка наличия куки
func (h *UserHandler) checkCookie(w http.ResponseWriter, r *http.Request) *models.User {
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

// isValidRegistrationFields - валидация полей регистрации
func isValidRegistrationFields(user *models.User) error {
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return errors.New("missing required fields")
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
	if user.Email == "" || user.Password == "" {
		return errors.New("missing required fields")
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

// RegisterUser - обработчик регистрации пользователя
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	userFromCoockies := h.checkCookie(w, r)
	if userFromCoockies != nil {
		log.Print("user already registered in")
		setCookie(w, userFromCoockies.Id)
		response.SendOKResponse(w)
		return
	}

	if r.Method != http.MethodPost {
		log.Print("method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print(err)
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w)
		return
	}

	err = isValidRegistrationFields(&user)
	if err != nil {
		log.Print(err)
		response.SendErrorResponse(err.Error(), http.StatusBadRequest, w)
		return
	}

	err = h.useCase.RegisterUser(&user)
	if err != nil {
		log.Print(err)
		response.SendErrorResponse(err.Error(), http.StatusBadRequest, w)
		return
	}

	setCookie(w, user.Id)
	response.SendOKResponse(w)
}

// LoginUser - обработчик авторизации пользователя
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	userFromCoockies := h.checkCookie(w, r)
	if userFromCoockies != nil {
		log.Print("user already logged in")
		setCookie(w, userFromCoockies.Id)
		response.SendOKResponse(w)
		return
	}

	if r.Method != http.MethodPost {
		log.Print("method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Print(err)
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w)
		return
	}

	err = isValidLoginFields(&user)
	if err != nil {
		log.Print(err)
		response.SendErrorResponse(err.Error(), http.StatusBadRequest, w)
		return
	}

	userId, err := h.useCase.AuthenticateUser(&user)
	if err != nil {
		log.Print(err)
		response.SendErrorResponse(err.Error(), http.StatusNotFound, w)
		return
	}

	setCookie(w, userId)
	response.SendOKResponse(w)
}

// LogoutUser - обработчик для выхода из сессии у пользователя
func (h *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	userFromCoockies := h.checkCookie(w, r)
	if userFromCoockies != nil {
		log.Printf("logout user %+v", userFromCoockies)
		deleteCookie(w, userFromCoockies.Id)

	}
	response.SendOKResponse(w)
}
