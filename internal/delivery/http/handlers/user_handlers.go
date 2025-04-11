package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"skillForce/internal/delivery/http/response"
	"skillForce/internal/models"
	"skillForce/internal/models/dto"
	"skillForce/pkg/logs"

	"github.com/badoux/checkmail"
)

// isValidRegistrationFields - валидация полей регистрации
func isValidRegistrationFields(user *dto.UserDTO) error {
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
func isValidLoginFields(user *dto.UserDTO) error {
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
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.PrintLog(r.Context(), "RegisterUser", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	var userInput dto.UserDTO
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		logs.PrintLog(r.Context(), "RegisterUser", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	err = isValidRegistrationFields(&userInput)
	if err != nil {
		logs.PrintLog(r.Context(), "RegisterUser", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusBadRequest, w, r)
		return
	}

	user := models.NewUser(userInput)
	cookie, err := h.useCase.RegisterUser(r.Context(), user)
	if err != nil {
		logs.PrintLog(r.Context(), "RegisterUser", fmt.Sprintf("%+v", err))

		if err.Error() == "email exists" { //TODO: хорошо бы все константы вывести в отдельный файл
			response.SendErrorResponse(err.Error(), http.StatusNotFound, w, r)
			return
		}

		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}

	logs.PrintLog(r.Context(), "RegisterUser", fmt.Sprintf("user %+v registered, send him cookie", user))

	setCookie(w, cookie)
	response.SendOKResponse(w, r)
}

// LoginUser godoc
// @Summary Login user
// @Description Login user with the given email, password and send cookie
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.UserDTO true "User information"
// @Success 200 {string} string "200 OK"
// @Failure 400 {object} response.ErrorResponse "invalid request | password too short | invalid email"
// @Failure 404 {object} response.ErrorResponse "email or password incorrect"
// @Failure 405 {object} response.ErrorResponse "method not allowed"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /api/login [post]
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.PrintLog(r.Context(), "LoginUser", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	var userInput dto.UserDTO
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		logs.PrintLog(r.Context(), "LoginUser", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	err = isValidLoginFields(&userInput)
	if err != nil {
		logs.PrintLog(r.Context(), "LoginUser", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusBadRequest, w, r)
		return
	}

	user := models.NewUser(userInput)
	cookieValue, err := h.useCase.AuthenticateUser(r.Context(), user)
	if err != nil {
		logs.PrintLog(r.Context(), "LoginUser", fmt.Sprintf("%+v", err))

		if err.Error() == "email or password incorrect" { //TODO: хорошо бы все константы вывести в отдельный файл
			response.SendErrorResponse(err.Error(), http.StatusNotFound, w, r)
			return
		}

		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}
	logs.PrintLog(r.Context(), "LoginUser", fmt.Sprintf("user %+v login, send him cookie", user))

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
func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	userProfile := h.checkCookie(r)
	if userProfile != nil {
		logs.PrintLog(r.Context(), "LogoutUser", fmt.Sprintf("logout user %+v", userProfile))
		err := h.useCase.LogoutUser(r.Context(), userProfile.Id)
		if err != nil {
			logs.PrintLog(r.Context(), "LogoutUser", fmt.Sprintf("%+v", err))
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
func (h *Handler) IsAuthorized(w http.ResponseWriter, r *http.Request) {
	userProfile := h.checkCookie(r)
	if userProfile != nil {
		userProfileOut := dto.UserProfileDTO{
			Name:      userProfile.Name,
			Email:     userProfile.Email,
			Bio:       userProfile.Bio,
			AvatarSrc: userProfile.AvatarSrc,
			HideEmail: userProfile.HideEmail,
		}

		logs.PrintLog(r.Context(), "IsAuthorized", fmt.Sprintf("user %+v is authorized", userProfile))
		response.SendUserProfile(&userProfileOut, w, r)
		return
	}
	logs.PrintLog(r.Context(), "IsAuthorized", "user not logged in")
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
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userProfile := h.checkCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "UpdateProfile", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	var UserProfileInput dto.UserProfileDTO
	err := json.NewDecoder(r.Body).Decode(&UserProfileInput)
	if err != nil {
		logs.PrintLog(r.Context(), "UpdateProfile", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	newUserProfile := models.NewUserProfile(UserProfileInput)
	//TODO: добавить тут валидацию
	err = h.useCase.UpdateProfile(r.Context(), userProfile.Id, newUserProfile)
	if err != nil {
		logs.PrintLog(r.Context(), "UpdateProfile", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}
	logs.PrintLog(r.Context(), "UpdateProfile", fmt.Sprintf("user %+v updated profile with values %+v", userProfile, newUserProfile))

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
// @Router /api/updateProfile [post]
func (h *Handler) UpdateProfilePhoto(w http.ResponseWriter, r *http.Request) {
	userProfile := h.checkCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "UpdateProfilePhoto", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		logs.PrintLog(r.Context(), "UpdateProfilePhoto", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("photo is too big", http.StatusBadRequest, w, r)
		return
	}

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		logs.PrintLog(r.Context(), "UpdateProfilePhoto", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("can`t reach photo", http.StatusBadRequest, w, r)
		return
	}
	defer file.Close()

	url, err := h.useCase.UploadFile(r.Context(), file, fileHeader)
	if err != nil {
		logs.PrintLog(r.Context(), "UpdateProfilePhoto", fmt.Sprintf("%+v", err))
		http.Error(w, "Ошибка загрузки в MinIO", http.StatusInternalServerError)
		return
	}

	newPhotoUrl, err := h.useCase.SaveProfilePhoto(r.Context(), url, userProfile.Id)
	if err != nil {
		logs.PrintLog(r.Context(), "UpdateProfilePhoto", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}

	logs.PrintLog(r.Context(), "UpdateProfilePhoto", fmt.Sprintf("Файл загружен: %s", newPhotoUrl))
	response.SendPhotoUrl(newPhotoUrl, w, r)
}

// DeleteProfilePhoto godoc
// @Summary Delete user`s profile photo
// @Description Deletes the profile photo of the authorized user and sets the default avatar
// @Tags users
// @Accept json
// @Produce json
// @Param profile body models.UserProfile true "Updated user profile"
// @Success 200 {string} string "200 OK"
// @Failure 400 {object} response.ErrorResponse "invalid request"
// @Failure 401 {object} response.ErrorResponse "not authorized"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /api/deleteProfilePhoto [post]
func (h *Handler) DeleteProfilePhoto(w http.ResponseWriter, r *http.Request) {
	userProfile := h.checkCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "DeleteProfilePhoto", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	err := h.useCase.DeleteProfilePhoto(r.Context(), userProfile.Id)
	if err != nil {
		logs.PrintLog(r.Context(), "DeleteProfilePhoto", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}

	logs.PrintLog(r.Context(), "DeleteProfilePhoto", fmt.Sprintf("Аватарка пользователя %+v заменена на стандартную", userProfile.Id))
	response.SendOKResponse(w, r)
}
