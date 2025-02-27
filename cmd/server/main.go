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
	userRepo := repository.NewUserRepository()
	userUseCase := usecase.NewUserUsecase(userRepo)
	userHandler := handlers.NewUserHandler(userUseCase)

	http.HandleFunc("/api/register", userHandler.RegisterUser)
	http.HandleFunc("/api/login", userHandler.LoginUser)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
