package mock

import (
	"fmt"
	"mime/multipart"
	"skillForce/backend/models"
)

type mockDB struct {
	users   map[int]*models.User
	courses map[int]*models.Course
	is_auth bool
}

func NewmockDB(is_auth bool) *mockDB {
	mockDB := &mockDB{
		users:   make(map[int]*models.User),
		courses: make(map[int]*models.Course),
		is_auth: is_auth,
	}
	for i := 1; i <= 4; i++ {
		course := models.Course{
			Id:              i,
			Price:           i,
			PurchasesAmount: i,
			CreatorId:       i,
			TimeToPass:      i,
			Title:           fmt.Sprintf("Курс #%d", i),
			Description:     fmt.Sprintf("Описание курса #%d", i),
			ScrImage:        fmt.Sprintf("image_%d.jpg", i),
		}
		mockDB.courses[i] = &course
	}

	for i := 1; i <= 4; i++ {
		user := models.User{
			Id:       i,
			Email:    fmt.Sprintf("user%d@skillforce.com", i),
			Password: fmt.Sprintf("password%d", i),
			Salt:     []byte(fmt.Sprintf("salt%d", i)),
		}
		mockDB.users[i] = &user
	}

	return mockDB
}

func (r *mockDB) GetUserByCookie(cookieValue string) (*models.UserProfile, error) {
	if r.is_auth {
		user := r.users[1]
		var userProfile models.UserProfile
		userProfile.Id = user.Id
		userProfile.Email = user.Email
		userProfile.Name = user.Name
		return &userProfile, nil
	}
	return nil, nil
}

func (r *mockDB) RegisterUser(user *models.User) (string, error) {
	return "", nil
}

func (r *mockDB) AuthenticateUser(email string, password string) (string, error) {
	return "", nil
}

func (r *mockDB) LogoutUser(userId int) error {
	return nil
}

func (r *mockDB) GetBucketCourses() ([]*models.Course, error) {
	var bucketCourses []*models.Course
	for i := 1; i <= 2; i++ {
		existingCourse := r.courses[i]
		bucketCourses = append(bucketCourses, existingCourse)
	}
	return bucketCourses, nil
}

func (r *mockDB) UpdateProfile(userId int, userProfile *models.UserProfile) error {
	return nil
}

func (r *mockDB) UpdateProfilePhoto(url string, userId int) error {
	return nil
}
func (r *mockDB) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	return "", nil
}
