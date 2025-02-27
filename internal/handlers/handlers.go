package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"skillForce/internal/models"
	"skillForce/internal/response"
	"skillForce/internal/usecase"
	"strconv"
	"time"
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

// RegisterUser - обработчик регистрации пользователя
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Вопрос: как лучше поступать, елси пользователь авторизован через куки?
	if r.Method != http.MethodPost {
		log.Print("method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	//TODO: реализовать валидацию
	if err != nil {
		log.Print(err)
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w)
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

	userId, err := h.useCase.AuthenticateUser(&user)
	if err != nil {
		log.Print(err)
		response.SendErrorResponse(err.Error(), http.StatusNotFound, w)
		return
	}

	setCookie(w, userId)
	response.SendOKResponse(w)
}
