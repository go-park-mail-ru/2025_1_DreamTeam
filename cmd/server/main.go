package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"skillForce/internal/handlers"
	"skillForce/internal/repository"
	"skillForce/internal/usecase"

	"github.com/joho/godotenv"
)

// точка входа приложения
func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Print(err)
		log.Fatal("Ошибка загрузки .env файла")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
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

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
