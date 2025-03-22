package main

import (
	"fmt"
	"log"
	"net/http"
	"skillForce/internal/env"
	"skillForce/internal/handlers"
	"skillForce/internal/middleware"
	"skillForce/internal/repository"
	"skillForce/internal/tools"
	"skillForce/internal/usecase"

	_ "skillForce/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	env := env.NewEnvironment()

	err := tools.InitMinio(env.MINIO_ENDPOINT, env.MINIO_ACCESS_KEY, env.MINIO_SECRET_ACCESS_KEY, env.MINIO_USESSL, env.MINIO_BUCKET_NAME)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", env.DB_HOST, env.DB_PORT, env.DB_USER, env.DB_PASSWORD, env.DB_NAME)
	database, err := repository.NewDatabase(dsn)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	siteMux := http.NewServeMux()

	userUseCase := usecase.NewUserUsecase(database)
	userHandler := handlers.NewUserHandler(userUseCase)

	courseUseCase := usecase.NewCourseUsecase(database)
	courseHandler := handlers.NewCourseHandler(courseUseCase)

	siteMux.HandleFunc("/api/register", userHandler.RegisterUser)
	siteMux.HandleFunc("/api/login", userHandler.LoginUser)
	siteMux.HandleFunc("/api/logout", userHandler.LogoutUser)
	siteMux.HandleFunc("/api/isAuthorized", userHandler.IsAuthorized)
	siteMux.HandleFunc("/api/updateProfile", userHandler.UpdateProfile)
	siteMux.HandleFunc("/api/updateProfilePhoto", userHandler.UpdateProfilePhoto)

	siteMux.HandleFunc("/api/getCourses", courseHandler.GetCourses)

	siteMux.HandleFunc("/api/docs/", httpSwagger.WrapHandler)

	//siteMux.HandleFunc("/api/panic", func(w http.ResponseWriter, r *http.Request) {
	//panic("this must me recovered")
	//})

	siteHandler := middleware.PanicMiddleware(siteMux)
	siteHandler = middleware.CorsOptionsMiddleware(siteHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", siteHandler))
}
