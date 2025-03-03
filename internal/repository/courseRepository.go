package repository

import (
	"fmt"
	"log"
	"math/rand"
	"skillForce/internal/models"
)

// CourseRepository - структура хранилища курсов
type CourseRepository struct {
	courses map[int]*models.Course
}

// NewCourseRepository - создание нового репозитория курсов
func NewCourseRepository() *CourseRepository {
	courses := make(map[int]*models.Course, 20)
	//TODO: убрать временную рыбу
	for i := 1; i <= 20; i++ {
		course := models.Course{
			Id:              i,
			Price:           rand.Intn(10000) + 1000, // Цена от 1000 до 11000
			PurchasesAmount: rand.Intn(500),          // Количество покупок от 0 до 500
			CreatorId:       rand.Intn(100) + 1,      // ID создателя от 1 до 100
			TimeToPass:      rand.Intn(50) + 5,       // Время прохождения от 5 до 55 часов
			Title:           fmt.Sprintf("Курс #%d", i),
			Description:     fmt.Sprintf("Описание курса #%d", i),
			ScrImage:        fmt.Sprintf("image_%d.jpg", i),
		}
		courses[i] = &course
	}

	return &CourseRepository{courses: courses}

}

// GetBucketCourses - извлекает список курсов из базы данных
func (r *CourseRepository) GetBucketCourses() ([]*models.Course, error) {
	var bucketCourses []*models.Course
	for _, existingCourse := range r.courses {
		bucketCourses = append(bucketCourses, existingCourse)
		if len(bucketCourses) == 16 {
			break
		}
	}
	log.Printf("get a bucket of courses: %+v", bucketCourses)
	return bucketCourses, nil
}
