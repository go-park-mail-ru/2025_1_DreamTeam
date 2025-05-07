package main

import (
	"log"
	"net/http"
	"skillForce/config"
	"skillForce/internal/delivery/http/middleware"
	"skillForce/pkg/logs"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	cookie "skillForce/internal/delivery/http/cookie"
	courseHandler "skillForce/internal/delivery/http/handlers/course"
	userHandler "skillForce/internal/delivery/http/handlers/user"

	courseUsecase "skillForce/internal/usecase/course"
	userUsecase "skillForce/internal/usecase/user"

	courseInfrastructure "skillForce/internal/repository/course"
	userInfrastructure "skillForce/internal/repository/user"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "skillForce/docs"
)

func main() {
	config := config.LoadConfig()

	siteMux := http.NewServeMux()

	userIfrustructure := userInfrastructure.NewUserInfrastructure(config)
	defer userIfrustructure.Close()
	userUsecase := userUsecase.NewUserUsecase(userIfrustructure)
	cookieManager := cookie.NewCookieManager(userUsecase)

	courseInfrastructure := courseInfrastructure.NewCourseInfrastructure(config)
	defer courseInfrastructure.Close()
	courseUsecase := courseUsecase.NewCourseUsecase(courseInfrastructure)
	courseHandler := courseHandler.NewHandler(cookieManager, courseUsecase)

	userHandler := userHandler.NewHandler(cookieManager)

	siteMux.HandleFunc("/api/register", userHandler.RegisterUser)
	siteMux.HandleFunc("/api/login", userHandler.LoginUser)
	siteMux.HandleFunc("/api/logout", userHandler.LogoutUser)
	siteMux.HandleFunc("/api/isAuthorized", userHandler.IsAuthorized)
	siteMux.Handle("/api/updateProfile", middleware.CSRFMiddleware(http.HandlerFunc(userHandler.UpdateProfile)))
	siteMux.Handle("/api/updateProfilePhoto", middleware.CSRFMiddleware(http.HandlerFunc(userHandler.UpdateProfilePhoto)))
	siteMux.Handle("/api/deleteProfilePhoto", middleware.CSRFMiddleware(http.HandlerFunc(userHandler.DeleteProfilePhoto)))
	siteMux.HandleFunc("/api/validEmail", userHandler.ConfirmUserEmail)

	siteMux.HandleFunc("/api/getCourses", courseHandler.GetCourses)
	siteMux.HandleFunc("/api/searchCourses", courseHandler.SearchCourses)
	siteMux.HandleFunc("/api/getCourse", courseHandler.GetCourse)
	siteMux.HandleFunc("/api/getCourseLesson", courseHandler.GetCourseLesson)
	siteMux.HandleFunc("/api/getNextLesson", courseHandler.GetNextLesson)
	siteMux.HandleFunc("/api/markLessonAsNotCompleted", courseHandler.MarkLessonAsNotCompleted)
	siteMux.HandleFunc("/api/getCourseRoadmap", courseHandler.GetCourseRoadmap)
	siteMux.HandleFunc("/api/video", courseHandler.ServeVideo)
	siteMux.HandleFunc("/api/createCourse", courseHandler.CreateCourse)
	siteMux.HandleFunc("/api/addCourseToFavourites", courseHandler.AddCourseToFavourites)
	siteMux.HandleFunc("/api/deleteCourseFromFavourites", courseHandler.DeleteCourseFromFavourites)
	siteMux.HandleFunc("/api/getFavouriteCourses", courseHandler.GetFavouriteCourses)
	siteMux.HandleFunc("/api/GetTestLesson", courseHandler.GetTestLesson)
	siteMux.HandleFunc("/api/AnswerQuiz", courseHandler.AnswerQuiz)
	siteMux.HandleFunc("/api/GetQuestionTestLesson", courseHandler.GetQuestionTestLesson)
	siteMux.HandleFunc("/api/AnswerQuestion", courseHandler.AnswerQuestion)

	siteMux.HandleFunc("/api/docs/", httpSwagger.WrapHandler)

	siteMux.Handle("/metrics", promhttp.Handler())

	siteHandler := logs.LoggerMiddleware(siteMux)
	siteHandler = middleware.MetricsMiddleware(siteHandler)
	siteHandler = middleware.PanicMiddleware(siteHandler)
	siteHandler = middleware.CorsOptionsMiddleware(siteHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", siteHandler))
}
