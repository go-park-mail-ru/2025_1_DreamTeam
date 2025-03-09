package main

import (
	"fmt"
	"log"
	"net/http"
	"skillForce/internal/env"
	"skillForce/internal/handlers"
	"skillForce/internal/repository"
	"skillForce/internal/usecase"
)

// точка входа приложения
func main() {
	env := env.NewEnvironment()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", env.DB_HOST, env.DB_PORT, env.DB_USER, env.DB_PASSWORD, env.DB_NAME)
	database, err := repository.NewDatabase(dsn)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	userUseCase := usecase.NewUserUsecase(database)
	userHandler := handlers.NewUserHandler(userUseCase)

	courseUseCase := usecase.NewCourseUsecase(database)
	courseHandler := handlers.NewCourseHandler(courseUseCase)

	http.HandleFunc("/api/register", userHandler.RegisterUser)
	http.HandleFunc("/api/login", userHandler.LoginUser)
	http.HandleFunc("/api/logout", userHandler.LogoutUser)

	http.HandleFunc("/api/getCourses", courseHandler.GetCourses)

	http.HandleFunc("/api/isAuthorized", userHandler.IsAuthorized)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
