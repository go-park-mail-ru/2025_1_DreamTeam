package main

import (
	"log"
	"net/http"
	"skillForce/config"
	"skillForce/internal/delivery/http/middleware"
	"skillForce/internal/repository/infrastructure"
	"skillForce/pkg/logs"

	cookie "skillForce/internal/delivery/http/cookie"
	courseHandler "skillForce/internal/delivery/http/handlers/course"
	userHandler "skillForce/internal/delivery/http/handlers/user"
	courseUsecase "skillForce/internal/usecase/course"
	userUsecase "skillForce/internal/usecase/user"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "skillForce/docs"
)

func main() {
	config := config.LoadConfig()

	infrastructure := infrastructure.NewInfrastructure(config)
	defer infrastructure.Close()

	siteMux := http.NewServeMux()

	userUsecase := userUsecase.NewUserUsecase(infrastructure)
	cookieManager := cookie.NewCookieManager(userUsecase)
	userHandler := userHandler.NewHandler(userUsecase, cookieManager)

	siteMux.HandleFunc("/api/register", userHandler.RegisterUser)
	siteMux.HandleFunc("/api/login", userHandler.LoginUser)
	siteMux.HandleFunc("/api/logout", userHandler.LogoutUser)
	siteMux.HandleFunc("/api/isAuthorized", userHandler.IsAuthorized)
	siteMux.HandleFunc("/api/updateProfile", userHandler.UpdateProfile)
	siteMux.HandleFunc("/api/updateProfilePhoto", userHandler.UpdateProfilePhoto)
	siteMux.HandleFunc("/api/deleteProfilePhoto", userHandler.DeleteProfilePhoto)
	siteMux.HandleFunc("/api/validEmail", userHandler.ConfirmUserEmail)

	courseUsecase := courseUsecase.NewCourseUsecase(infrastructure)
	courseHandler := courseHandler.NewHandler(courseUsecase, cookieManager)

	siteMux.HandleFunc("/api/getCourses", courseHandler.GetCourses)
	siteMux.HandleFunc("/api/getCourse", courseHandler.GetCourse)
	siteMux.HandleFunc("/api/getCourseLesson", courseHandler.GetCourseLesson)
	siteMux.HandleFunc("/api/getNextLesson", courseHandler.GetNextLesson)
	siteMux.HandleFunc("/api/markLessonAsNotCompleted", courseHandler.MarkLessonAsNotCompleted)
	siteMux.HandleFunc("/api/getCourseRoadmap", courseHandler.GetCourseRoadmap)
	siteMux.HandleFunc("/api/video", courseHandler.ServeVideo)

	siteMux.HandleFunc("/api/docs/", httpSwagger.WrapHandler)

	siteHandler := logs.LoggerMiddleware(siteMux)
	siteHandler = middleware.PanicMiddleware(siteHandler)
	siteHandler = middleware.CorsOptionsMiddleware(siteHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", siteHandler))
}
