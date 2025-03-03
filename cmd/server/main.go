package main

import (
	"log"
	"net/http"
	"skillForce/internal/handlers"
	"skillForce/internal/repository"
	"skillForce/internal/usecase"
)

// точка входа приложения
func main() {
	// userRepo := repository.NewUserRepository()
	// userUseCase := usecase.NewUserUsecase(userRepo)
	// userHandler := handlers.NewUserHandler(userUseCase)

	database, err := repository.NewDatabase("host=localhost port=5432 user=dmitrii password=password dbname=skillforce_test")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	database.GetUserByCookie("")

	userUseCase := usecase.NewUserUsecase(database)
	userHandler := handlers.NewUserHandler(userUseCase)
	// courseRepo := repository.NewCourseRepository()
	courseUseCase := usecase.NewCourseUsecase(database)
	courseHandler := handlers.NewCourseHandler(courseUseCase)

	http.HandleFunc("/api/register", userHandler.RegisterUser)
	http.HandleFunc("/api/login", userHandler.LoginUser)
	http.HandleFunc("/api/logout", userHandler.LogoutUser)

	http.HandleFunc("/api/getCourses", courseHandler.GetCourses)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
