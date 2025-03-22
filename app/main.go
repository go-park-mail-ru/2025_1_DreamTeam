package main

import (
	"log"
	"net/http"
	"skillForce/backend/delivery/http/middleware"
	"skillForce/backend/repository/infrastructure"
	"skillForce/backend/usecase"
	"skillForce/env"

	_ "skillForce/docs"

	"skillForce/backend/delivery/http/handlers"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	env := env.NewEnvironment()

	infrastructure := infrastructure.NewInfrastructure(env)
	defer infrastructure.Close()

	siteMux := http.NewServeMux()

	userUseCase := usecase.NewUserUsecase(infrastructure)
	userHandler := handlers.NewUserHandler(userUseCase)

	courseUseCase := usecase.NewCourseUsecase(infrastructure)
	courseHandler := handlers.NewCourseHandler(courseUseCase)

	siteMux.HandleFunc("/api/register", userHandler.RegisterUser)
	siteMux.HandleFunc("/api/login", userHandler.LoginUser)
	siteMux.HandleFunc("/api/logout", userHandler.LogoutUser)
	siteMux.HandleFunc("/api/isAuthorized", userHandler.IsAuthorized)
	siteMux.HandleFunc("/api/updateProfile", userHandler.UpdateProfile)
	siteMux.HandleFunc("/api/updateProfilePhoto", userHandler.UpdateProfilePhoto)

	siteMux.HandleFunc("/api/getCourses", courseHandler.GetCourses)

	siteMux.HandleFunc("/api/docs/", httpSwagger.WrapHandler)

	siteHandler := middleware.PanicMiddleware(siteMux)
	siteHandler = middleware.CorsOptionsMiddleware(siteHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", siteHandler))
}
