package main

import (
	"fmt"
	"log"
	"net/http"
	"skillForce/internal/env"
	"skillForce/internal/handlers"
	"skillForce/internal/middleware"
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

	corsMux := http.NewServeMux()

	userUseCase := usecase.NewUserUsecase(database)
	userHandler := handlers.NewUserHandler(userUseCase)

	courseUseCase := usecase.NewCourseUsecase(database)
	courseHandler := handlers.NewCourseHandler(courseUseCase)

	corsMux.HandleFunc("/api/register", userHandler.RegisterUser)
	corsMux.HandleFunc("/api/login", userHandler.LoginUser)
	corsMux.HandleFunc("/api/logout", userHandler.LogoutUser)
	corsMux.HandleFunc("/api/isAuthorized", userHandler.IsAuthorized)

	corsMux.HandleFunc("/api/getCourses", courseHandler.GetCourses)

	corsHandler := middleware.CorsOptionsMiddleware(corsMux)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
