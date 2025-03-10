package repository

import (
	"fmt"
	"skillForce/internal/models"
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

func (r *mockDB) GetUserByCookie(cookieValue string) (*models.User, error) {
	if r.is_auth {
		user := r.users[1]
		return user, nil
	}
	return nil, nil
}

func (r *mockDB) RegisterUser(user *models.User) error {
	return nil
}

func (r *mockDB) AuthenticateUser(email string, password string) (int, error) {
	return 0, nil
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
