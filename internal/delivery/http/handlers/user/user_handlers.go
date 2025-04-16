package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"skillForce/internal/delivery/http/response"
	"skillForce/internal/models/dto"
	models "skillForce/internal/models/user"
	"skillForce/pkg/logs"

	"github.com/badoux/checkmail"
)

type UserUsecaseInterface interface {
	RegisterUser(ctx context.Context, token string) (string, error)
	ValidUser(ctx context.Context, user *models.User) error
	AuthenticateUser(ctx context.Context, user *models.User) (string, error)
	GetUserByCookie(ctx context.Context, cookieValue string) (*models.UserProfile, error)
	LogoutUser(ctx context.Context, userId int) error
	UpdateProfile(ctx context.Context, userId int, userProfile *models.UserProfile) error
	UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	SaveProfilePhoto(ctx context.Context, url string, userId int) (string, error)
	DeleteProfilePhoto(ctx context.Context, userId int) error
}

type CookieManagerInterface interface {
	CheckCookie(r *http.Request) *models.UserProfile
	SetCookie(w http.ResponseWriter, cookieValue string)
	DeleteCookie(w http.ResponseWriter)
}

type Handler struct {
	userUsecase   UserUsecaseInterface
	cookieManager CookieManagerInterface
}

func NewHandler(userUsecase UserUsecaseInterface, cookieManager CookieManagerInterface) *Handler {
	return &Handler{
		userUsecase:   userUsecase,
		cookieManager: cookieManager,
	}
}

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

	user := &models.User{
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: userInput.Password,
	}
	err = h.userUsecase.ValidUser(r.Context(), user)
	if err != nil {
		if errors.Is(err, errors.New("email exists")) {
			logs.PrintLog(r.Context(), "RegisterUser", fmt.Sprintf("%+v", err))
			response.SendErrorResponse("email exists", http.StatusNotFound, w, r)
			return
		}

		logs.PrintLog(r.Context(), "RegisterUser", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusConflict, w, r)
		return
	}

	logs.PrintLog(r.Context(), "RegisterUser", fmt.Sprintf("send email to confirm user %+v", user))
	response.SendOKResponse(w, r)
}

// ConfirmUserEmail godoc
// @Summary Confirm user email
// @Description Confirm user email using the token from the registration email
// @Tags users
// @Accept json
// @Produce json
// @Param token query string true "Token from registration email"
// @Success 200 {string} string "200 OK"
// @Failure 400 {object} response.ErrorResponse "invalid token"
// @Failure 404 {object} response.ErrorResponse "email exists"
// @Failure 405 {object} response.ErrorResponse "method not allowed"
// @Failure 500 {object} response.ErrorResponse "server error"
// @Router /api/validEmail [get]
func (h *Handler) ConfirmUserEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logs.PrintLog(r.Context(), "ConfirmUserEmail", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		logs.PrintLog(r.Context(), "ConfirmUserEmail", "invalid token")
		response.SendErrorResponse("invalid token", http.StatusBadRequest, w, r)
		return
	}

	cookie, err := h.userUsecase.RegisterUser(r.Context(), token)
	if err != nil {
		if errors.Is(err, errors.New("invalid token")) {
			logs.PrintLog(r.Context(), "ConfirmUserEmail", fmt.Sprintf("%+v", err))
			response.SendErrorResponse("invalid token", http.StatusBadRequest, w, r)
			return
		}
		if errors.Is(err, errors.New("email exists")) {
			logs.PrintLog(r.Context(), "ConfirmUserEmail", fmt.Sprintf("%+v", err))
			response.SendErrorResponse("email exists", http.StatusNotFound, w, r)
			return
		}
		logs.PrintLog(r.Context(), "ConfirmUserEmail", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}
	logs.PrintLog(r.Context(), "ConfirmUserEmail", "register user and send him cookie")
	h.cookieManager.SetCookie(w, cookie)
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

	user := &models.User{
		Email:    userInput.Email,
		Password: userInput.Password,
	}
	cookieValue, err := h.userUsecase.AuthenticateUser(r.Context(), user)
	if err != nil {
		logs.PrintLog(r.Context(), "LoginUser", fmt.Sprintf("%+v", err))

		if errors.Is(err, errors.New("email or password incorrect")) {
			response.SendErrorResponse(err.Error(), http.StatusNotFound, w, r)
			return
		}

		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}
	logs.PrintLog(r.Context(), "LoginUser", fmt.Sprintf("user %+v login, send him cookie", user))

	h.cookieManager.SetCookie(w, cookieValue)
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
	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile != nil {
		logs.PrintLog(r.Context(), "LogoutUser", fmt.Sprintf("logout user %+v", userProfile))
		err := h.userUsecase.LogoutUser(r.Context(), userProfile.Id)
		if err != nil {
			logs.PrintLog(r.Context(), "LogoutUser", fmt.Sprintf("%+v", err))
			response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
			return
		}
		h.cookieManager.DeleteCookie(w)

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
	userProfile := h.cookieManager.CheckCookie(r)
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
	if r.Method == http.MethodGet {
		response.SendOKResponse(w, r)
		return
	}
	userProfile := h.cookieManager.CheckCookie(r)
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

	newUserProfile := &models.UserProfile{
		Name:      UserProfileInput.Name,
		Bio:       UserProfileInput.Bio,
		Email:     UserProfileInput.Email,
		AvatarSrc: UserProfileInput.AvatarSrc,
		HideEmail: UserProfileInput.HideEmail,
	}
	//TODO: добавить тут валидацию
	err = h.userUsecase.UpdateProfile(r.Context(), userProfile.Id, newUserProfile)
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
	if r.Method == http.MethodGet {
		response.SendOKResponse(w, r)
		return
	}
	userProfile := h.cookieManager.CheckCookie(r)
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

	url, err := h.userUsecase.UploadFile(r.Context(), file, fileHeader)
	if err != nil {
		logs.PrintLog(r.Context(), "UpdateProfilePhoto", fmt.Sprintf("%+v", err))
		http.Error(w, "Ошибка загрузки в MinIO", http.StatusInternalServerError)
		return
	}

	newPhotoUrl, err := h.userUsecase.SaveProfilePhoto(r.Context(), url, userProfile.Id)
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
	if r.Method == http.MethodGet {
		response.SendOKResponse(w, r)
		return
	}
	userProfile := h.cookieManager.CheckCookie(r)
	if userProfile == nil {
		logs.PrintLog(r.Context(), "DeleteProfilePhoto", "user not logged in")
		response.SendErrorResponse("not authorized", http.StatusUnauthorized, w, r)
		return
	}

	err := h.userUsecase.DeleteProfilePhoto(r.Context(), userProfile.Id)
	if err != nil {
		logs.PrintLog(r.Context(), "DeleteProfilePhoto", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}

	logs.PrintLog(r.Context(), "DeleteProfilePhoto", fmt.Sprintf("Аватарка пользователя %+v заменена на стандартную", userProfile.Id))
	response.SendOKResponse(w, r)
}
