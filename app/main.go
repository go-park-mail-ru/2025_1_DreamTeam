package main

import (
	"log"
	"net/http"
	"skillForce/config"
	"skillForce/internal/delivery/http/middleware"
	"skillForce/internal/repository/infrastructure"
	"skillForce/internal/usecase"
	"skillForce/pkg/logs"

	_ "skillForce/docs"

	"skillForce/internal/delivery/http/handlers"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	config := config.LoadConfig()

	infrastructure := infrastructure.NewInfrastructure(config)
	defer infrastructure.Close()

	siteMux := http.NewServeMux()

	useCase := usecase.NewUsecase(infrastructure)
	handler := handlers.NewHandler(useCase)

	siteMux.HandleFunc("/api/register", handler.RegisterUser)
	siteMux.HandleFunc("/api/login", handler.LoginUser)
	siteMux.HandleFunc("/api/logout", handler.LogoutUser)
	siteMux.HandleFunc("/api/isAuthorized", handler.IsAuthorized)
	siteMux.HandleFunc("/api/updateProfile", handler.UpdateProfile)
	siteMux.HandleFunc("/api/updateProfilePhoto", handler.UpdateProfilePhoto)
	siteMux.HandleFunc("/api/deleteProfilePhoto", handler.DeleteProfilePhoto)

	siteMux.HandleFunc("/api/getCourses", handler.GetCourses)
	siteMux.HandleFunc("/api/getCourse", handler.GetCourse)
	siteMux.HandleFunc("/api/getCourseLesson", handler.GetCourseLesson)
	siteMux.HandleFunc("/api/getNextLesson", handler.GetNextLesson)
	siteMux.HandleFunc("/api/markLessonAsNotCompleted", handler.MarkLessonAsNotCompleted)
	siteMux.HandleFunc("/api/getCourseRoadmap", handler.GetCourseRoadmap)
	siteMux.HandleFunc("/api/video", handler.ServeVideo)

	siteMux.HandleFunc("/api/docs/", httpSwagger.WrapHandler)

	siteHandler := logs.LoggerMiddleware(siteMux)
	siteHandler = middleware.PanicMiddleware(siteHandler)
	siteHandler = middleware.CorsOptionsMiddleware(siteHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", siteHandler))
}
