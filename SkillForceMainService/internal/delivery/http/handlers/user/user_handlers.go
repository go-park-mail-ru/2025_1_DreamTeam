package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	userpb "skillForce/internal/delivery/grpc/proto/user"
	"skillForce/internal/delivery/http/response"
	"skillForce/internal/models/dto"
	models "skillForce/internal/models/user"
	"skillForce/pkg/logs"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/mailru/easyjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type CookieManagerInterface interface {
	CheckCookie(r *http.Request) *models.UserProfile
	SetCookie(w http.ResponseWriter, cookieValue string)
	DeleteCookie(w http.ResponseWriter)
	LogoutUser(ctx context.Context, userId int) error
}

type Handler struct {
	userClient    userpb.UserServiceClient
	cookieManager CookieManagerInterface
}

func NewHandler(cookieManager CookieManagerInterface) *Handler {
	conn, err := grpc.NewClient("user-service:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	userClient := userpb.NewUserServiceClient(conn)

	return &Handler{
		cookieManager: cookieManager,
		userClient:    userClient,
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
// @Param user body usermodels.User true "User information"
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
	if err := easyjson.UnmarshalFromReader(r.Body, &userInput); err != nil {
		logs.PrintLog(r.Context(), "RegisterUser", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	err := isValidRegistrationFields(&userInput)
	if err != nil {
		logs.PrintLog(r.Context(), "RegisterUser", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusBadRequest, w, r)
		return
	}

	grpcUser := &userpb.User{
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: userInput.Password,
	}

	_, err = h.userClient.ValidUser(r.Context(), grpcUser)

	if err != nil {
		if strings.Contains(err.Error(), "email exists") {
			logs.PrintLog(r.Context(), "RegisterUser", "email exists")
			response.SendErrorResponse("email exists", http.StatusNotFound, w, r)
			return
		}
		logs.PrintLog(r.Context(), "RegisterUser", fmt.Sprintf("gRPC error: %+v", err))
		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}

	logs.PrintLog(r.Context(), "RegisterUser", fmt.Sprintf("send email to confirm user %+v", userInput.Email))
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

	grpcToken := &userpb.RegisterRequest{
		Token: token,
	}
	registerResp, err := h.userClient.RegisterUser(r.Context(), grpcToken)
	if err != nil {
		if strings.Contains(err.Error(), "invalid token") {
			logs.PrintLog(r.Context(), "ConfirmUserEmail", fmt.Sprintf("%+v", err))
			response.SendErrorResponse("invalid token", http.StatusBadRequest, w, r)
			return
		}
		if strings.Contains(err.Error(), "email exists") {
			logs.PrintLog(r.Context(), "ConfirmUserEmail", fmt.Sprintf("%+v", err))
			response.SendErrorResponse("email exists", http.StatusNotFound, w, r)
			return
		}
		logs.PrintLog(r.Context(), "ConfirmUserEmail", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}
	logs.PrintLog(r.Context(), "ConfirmUserEmail", "register user and send him cookie")
	h.cookieManager.SetCookie(w, registerResp.CookieVal)
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
	if err := easyjson.UnmarshalFromReader(r.Body, &userInput); err != nil {
		logs.PrintLog(r.Context(), "LoginUser", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	err := isValidLoginFields(&userInput)
	if err != nil {
		logs.PrintLog(r.Context(), "LoginUser", fmt.Sprintf("%+v", err))
		response.SendErrorResponse(err.Error(), http.StatusBadRequest, w, r)
		return
	}

	grpcUser := &userpb.User{
		Email:    userInput.Email,
		Password: userInput.Password,
	}
	authResp, err := h.userClient.AuthenticateUser(r.Context(), grpcUser)

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			logs.PrintLog(r.Context(), "LoginUser", fmt.Sprintf("Not grpc error: %+v", err))
			response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
			return
		}

		if st.Message() == "email or password incorrect" {
			logs.PrintLog(r.Context(), "LoginUser", fmt.Sprintf("%+v", err))
			response.SendErrorResponse("email or password incorrect", http.StatusNotFound, w, r)
			return
		}

		logs.PrintLog(r.Context(), "LoginUser", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}
	h.cookieManager.SetCookie(w, authResp.CookieVal)
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
		err := h.cookieManager.LogoutUser(r.Context(), userProfile.Id)
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
			IsAdmin:   userProfile.IsAdmin,
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
// @Param profile body usermodels.UserProfile true "Updated user profile"
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
	if err := easyjson.UnmarshalFromReader(r.Body, &UserProfileInput); err != nil {
		logs.PrintLog(r.Context(), "UpdateProfile", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("invalid request", http.StatusBadRequest, w, r)
		return
	}

	grpcNewUserProfile := &userpb.UserProfile{
		Name:      UserProfileInput.Name,
		Bio:       UserProfileInput.Bio,
		Email:     UserProfileInput.Email,
		AvatarSrc: UserProfileInput.AvatarSrc,
		HideEmail: UserProfileInput.HideEmail,
	}

	grpcUpdateProfileRequest := &userpb.UpdateProfileRequest{
		UserId:  int32(userProfile.Id),
		Profile: grpcNewUserProfile,
	}

	_, err := h.userClient.UpdateProfile(r.Context(), grpcUpdateProfileRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "UpdateProfile", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}

	newUserProfile := dto.UserProfileDTO{
		Name:      grpcNewUserProfile.Name,
		Bio:       grpcNewUserProfile.Bio,
		Email:     grpcNewUserProfile.Email,
		AvatarSrc: grpcNewUserProfile.AvatarSrc,
		HideEmail: grpcNewUserProfile.HideEmail,
		IsAdmin:   grpcNewUserProfile.IsAdmin,
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
// @Param profile body usermodels.UserProfile true "Updated user profile"
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
	defer func() {
		if err := file.Close(); err != nil {
			logs.PrintLog(r.Context(), "ServeVideo", "failed to close reader")
		}
	}()

	fileData, err := io.ReadAll(file)
	if err != nil {
		logs.PrintLog(r.Context(), "UpdateProfilePhoto", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("can`t reach photo", http.StatusBadRequest, w, r)
		return
	}
	grpcUploadFileRequest := &userpb.UploadFileRequest{
		FileName:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		FileData:    fileData,
	}

	grpcUploadFileResp, err := h.userClient.UploadFile(r.Context(), grpcUploadFileRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "UpdateProfilePhoto", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}

	grpcSaveProfilePhotoRequest := &userpb.SaveProfilePhotoRequest{
		UserId: int32(userProfile.Id),
		Url:    grpcUploadFileResp.Url,
	}

	grpcSaveProfilePhotoResp, err := h.userClient.SaveProfilePhoto(r.Context(), grpcSaveProfilePhotoRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "UpdateProfilePhoto", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}

	logs.PrintLog(r.Context(), "UpdateProfilePhoto", fmt.Sprintf("Файл загружен: %s", grpcSaveProfilePhotoResp.NewPhtotoUrl))
	response.SendPhotoUrl(grpcSaveProfilePhotoResp.NewPhtotoUrl, w, r)
}

// DeleteProfilePhoto godoc
// @Summary Delete user`s profile photo
// @Description Deletes the profile photo of the authorized user and sets the default avatar
// @Tags users
// @Accept json
// @Produce json
// @Param profile body usermodels.UserProfile true "Updated user profile"
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

	grpcDeleteProfilePhotoRequest := &userpb.DeleteProfilePhotoRequest{
		UserId: int32(userProfile.Id),
	}

	_, err := h.userClient.DeleteProfilePhoto(r.Context(), grpcDeleteProfilePhotoRequest)
	if err != nil {
		logs.PrintLog(r.Context(), "DeleteProfilePhoto", fmt.Sprintf("%+v", err))
		response.SendErrorResponse("server error", http.StatusInternalServerError, w, r)
		return
	}

	logs.PrintLog(r.Context(), "DeleteProfilePhoto", fmt.Sprintf("Аватарка пользователя %+v заменена на стандартную", userProfile.Id))
	response.SendOKResponse(w, r)
}
