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

// checkCookie - проверка наличия куки
func (h *UserHandler) checkCookie(w http.ResponseWriter, r *http.Request) bool {
	session, err := r.Cookie("session_id")
	loggedIn := (err != http.ErrNoCookie)
	if loggedIn {
		user, err := h.useCase.GetUserByCookie(session.Value)
		if err != nil {
			log.Print(err)
			return false
		}
		setCookie(w, user.Id)
		return true
	}
	return false
}

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
	if h.checkCookie(w, r) {
		log.Print("user already registered in")
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
	if h.checkCookie(w, r) {
		log.Print("user already logged in")
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
