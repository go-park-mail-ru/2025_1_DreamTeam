package usecase

import (
	"skillForce/internal/hash"
	"skillForce/internal/models"
	"skillForce/internal/repository"
)

// UserUsecase - структура бизнес-логики, которая ожидает интерфейс репозитория
type UserUsecase struct {
	repo repository.Repository
}

// NewUserUsecase - конструктор
func NewUserUsecase(repo repository.Repository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// RegisterUser - регистрация пользователя
func (uc *UserUsecase) RegisterUser(user *models.User) (string, error) {
	err := hash.HashPasswordAndCreateSalt(user)
	if err != nil {
		return "", err
	}

	return uc.repo.RegisterUser(user)
}

// AuthenticateUser - авторизация пользователя
func (uc *UserUsecase) AuthenticateUser(user *models.User) (string, error) {
	return uc.repo.AuthenticateUser(user.Email, user.Password)
}

// GetUserByCookie - получение пользователя по cookie
func (uc *UserUsecase) GetUserByCookie(cookieValue string) (*models.UserProfile, error) {
	return uc.repo.GetUserByCookie(cookieValue)
}

func (uc *UserUsecase) LogoutUser(userId int) error {
	return uc.repo.LogoutUser(userId)
}

func (uc *UserUsecase) UpdateProfile(userId int, userProfile *models.UserProfile) error {
	return uc.repo.UpdateProfile(userId, userProfile)
}

func (uc *UserUsecase) SaveProfilePhoto(url string, userId int) error {
	// minioClient, err :=
	// if err != nil {
	// 	return "", err
	// }

	// objectName := fileHeader.Filename // Можно добавить уникальность
	// contentType := fileHeader.Header.Get("Content-Type")

	// // Загрузка файла в MinIO
	// info, err := minioClient.PutObject(
	// 	context.Background(),
	// 	bucketName,
	// 	objectName,
	// 	file,
	// 	fileHeader.Size,
	// 	minio.PutObjectOptions{ContentType: contentType},
	// )
	// if err != nil {
	// 	return "", err
	// }

	// fileURL := fmt.Sprintf("https://%s/%s/%s", endpoint, bucketName, info.Key)
	// return fileURL, nil
	return uc.repo.UpdateProfilePhoto(url, userId)
}
