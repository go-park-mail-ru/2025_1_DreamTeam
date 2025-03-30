package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"skillForce/backend/delivery/http/response"
	"skillForce/backend/models"
	"skillForce/backend/usecase"
	"time"

	"github.com/badoux/checkmail"
)

// UserHandler - структура обработчика HTTP-запросов
type UserHandler struct {
	useCase usecase.UserUsecaseInterface
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

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user with the given name, email, and password and send cookie
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User information"
// @Success 200 {string} string "200 OK"
// @Failure 400 {object} response.ErrorResponse "invalid request | missing required fields | password too short | invalid email"
// @Failure 404 {object} response.ErrorResponse "email exists"
// @Failure 405 {object} response.ErrorResponse "method not allowed"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /api/register [post]
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

// LoginUser godoc
// @Summary Login user
// @Description Login user with the given email, password and send cookie
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User information"
// @Success 200 {string} string "200 OK"
// @Failure 400 {object} response.ErrorResponse "invalid request | password too short | invalid email"
// @Failure 404 {object} response.ErrorResponse "email or password incorrect"
// @Failure 405 {object} response.ErrorResponse "method not allowed"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /api/login [post]
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

		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}

	log.Printf("user %v login, send him cookie", user)

	setCookie(w, cookieValue)
	response.SendOKResponse(w, r)
}

// LogoutUser godoc
// @Summary Logout user
// @Description Logout user by deleting session cookie
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {string} string "200 OK"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /api/logout [post]
func (h *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	userProfile := h.checkCookie(r)
	if userProfile != nil {
		log.Printf("logout user %+v", userProfile)
		err := h.useCase.LogoutUser(userProfile.Id)
		if err != nil {
			log.Printf("from logoutUser: %v", err)
			response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
			return
		}
		deleteCookie(w)

	}
	response.SendOKResponse(w, r)
}

// IsAuthorized godoc
// @Summary Check if user is authorized
// @Description Returns user profile if authorized
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} response.UserProfileResponse "User profile"
// @Failure 401 {object} response.ErrorResponse "not authorized"
// @Router /api/isAuthorized [get]
func (h *UserHandler) IsAuthorized(w http.ResponseWriter, r *http.Request) {
	userProfile := h.checkCookie(r)
	if userProfile != nil {
		response.SendUserProfile(userProfile, w, r)
		return
	}
	log.Print("from isAuthorized: user not logged in")
	response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Updates the profile information of the authorized user
// @Tags users
// @Accept json
// @Produce json
// @Param profile body models.UserProfile true "Updated user profile"
// @Success 200 {string} string "200 OK"
// @Failure 400 {object} response.ErrorResponse "invalid request"
// @Failure 401 {object} response.ErrorResponse "not authorized"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /api/updateProfile [post]
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

// UpdateProfilePhoto godoc
// @Summary Update user profile photo
// @Description Updates the profile photo of the authorized user
// @Tags users
// @Accept json
// @Produce json
// @Param profile body models.UserProfile true "Updated user profile"
// @Success 200 {string} string "200 OK"
// @Failure 400 {object} response.ErrorResponse "invalid request"
// @Failure 401 {object} response.ErrorResponse "not authorized"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /api/updateProfilePhoto [post]
func (h *UserHandler) UpdateProfilePhoto(w http.ResponseWriter, r *http.Request) {
	userProfile := h.checkCookie(r)
	if userProfile == nil {
		log.Print("from updateProfile: user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		log.Printf("from updateProfilePhoto: %v", err)
		response.SendErrorResponse("photo is too big", http.StatusBadRequest, w, r)
		return
	}

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		log.Printf("from updateProfilePhoto: %v", err)
		response.SendErrorResponse("can`t reach photo", http.StatusBadRequest, w, r)
		return
	}
	defer file.Close()

	url, err := h.useCase.UploadFile(file, fileHeader)
	if err != nil {
		http.Error(w, "Ошибка загрузки в MinIO", http.StatusInternalServerError)
		return
	}

	err = h.useCase.SaveProfilePhoto(url, userProfile.Id)
	if err != nil {
		log.Printf("from updateProfilePhoto: %v", err)
		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}

	fmt.Fprintf(w, "Файл загружен: %s", url)
}
