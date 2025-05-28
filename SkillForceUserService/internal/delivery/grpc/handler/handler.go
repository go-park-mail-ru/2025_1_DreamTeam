package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	userpb "skillForce/internal/delivery/grpc/proto"
	models "skillForce/internal/models/user"
	"skillForce/pkg/logs"

	"google.golang.org/protobuf/types/known/emptypb"
)

type UserUsecaseInterface interface {
	RegisterUser(ctx context.Context, token string) (string, error)
	ValidUser(ctx context.Context, user *models.User) error
	AuthenticateUser(ctx context.Context, user *models.User) (string, error)
	LogoutUser(ctx context.Context, userId int) error
	UpdateProfile(ctx context.Context, userId int, userProfile *models.UserProfile) error
	UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	SaveProfilePhoto(ctx context.Context, url string, userId int) (string, error)
	DeleteProfilePhoto(ctx context.Context, userId int) error
}

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	usecase UserUsecaseInterface
}

func NewUserHandler(uc UserUsecaseInterface) *UserHandler {
	return &UserHandler{usecase: uc}
}

// RegisterUser handles user registration by token
func (h *UserHandler) RegisterUser(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	cookieVal, err := h.usecase.RegisterUser(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &userpb.RegisterResponse{CookieVal: cookieVal}, nil
}

// ValidUser handles email confirmation
func (h *UserHandler) ValidUser(ctx context.Context, req *userpb.User) (*emptypb.Empty, error) {
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	err := h.usecase.ValidUser(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// AuthenticateUser handles login
func (h *UserHandler) AuthenticateUser(ctx context.Context, req *userpb.User) (*userpb.AuthenticateResponse, error) {
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	cookieVal, err := h.usecase.AuthenticateUser(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &userpb.AuthenticateResponse{CookieVal: cookieVal}, nil
}

// UpdateProfile handles profile update
func (h *UserHandler) UpdateProfile(ctx context.Context, req *userpb.UpdateProfileRequest) (*emptypb.Empty, error) {
	userProfile := models.UserProfile{
		Id:        int(req.UserId),
		Email:     req.Profile.Email,
		Name:      req.Profile.Name,
		Bio:       req.Profile.Bio,
		HideEmail: req.Profile.HideEmail,
	}
	err := h.usecase.UpdateProfile(ctx, int(req.UserId), &userProfile)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// UploadFile handles file upload
func (h *UserHandler) UploadFile(ctx context.Context, req *userpb.UploadFileRequest) (*userpb.UploadFileResponse, error) {
	file, fileHeader, err := ConvertToMultipart(req.FileData, req.FileName, req.ContentType)
	if err != nil {
		return nil, err
	}
	url, err := h.usecase.UploadFile(ctx, file, fileHeader)
	if err != nil {
		return nil, err
	}
	return &userpb.UploadFileResponse{Url: url}, nil
}

// SaveProfilePhoto saves uploaded profile photo URL
func (h *UserHandler) SaveProfilePhoto(ctx context.Context, req *userpb.SaveProfilePhotoRequest) (*userpb.SaveProfilePhotoResponse, error) {
	newURL, err := h.usecase.SaveProfilePhoto(ctx, req.Url, int(req.UserId))
	if err != nil {
		return nil, err
	}
	return &userpb.SaveProfilePhotoResponse{NewPhtotoUrl: newURL}, nil
}

// DeleteProfilePhoto removes profile photo
func (h *UserHandler) DeleteProfilePhoto(ctx context.Context, req *userpb.DeleteProfilePhotoRequest) (*emptypb.Empty, error) {
	err := h.usecase.DeleteProfilePhoto(ctx, int(req.UserId))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func ConvertToMultipart(fileData []byte, fileName, contentType string) (multipart.File, *multipart.FileHeader, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Создаем заголовки для файла
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="`+fileName+`"`)
	h.Set("Content-Type", contentType)

	// Пишем файл в multipart.Writer
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, nil, err
	}
	if _, err := io.Copy(part, bytes.NewReader(fileData)); err != nil {
		return nil, nil, err
	}
	err = writer.Close()
	if err != nil {
		logs.PrintLog(context.Background(), "ConvertToMultipart", fmt.Sprintf("%+v", err))
	}

	// Парсим то, что получилось, как multipart/form-data
	r := multipart.NewReader(body, writer.Boundary())
	form, err := r.ReadForm(int64(len(fileData)) + 1024) // выделяем буфер
	if err != nil {
		return nil, nil, err
	}

	files := form.File["file"]
	if len(files) == 0 {
		return nil, nil, io.EOF
	}
	fh := files[0]
	f, err := fh.Open()
	if err != nil {
		return nil, nil, err
	}

	return f, fh, nil
}
