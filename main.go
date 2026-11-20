package main

import (
	"log"
	"os"

	"gc2-university-api/config"
	_ "gc2-university-api/docs"
	"gc2-university-api/handler"
	"gc2-university-api/middleware"
	"gc2-university-api/repository"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title University Academic Program API
// @version 1.0
// @description REST API for managing students, courses, and enrollments.
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name authorization
func main() {
	_ = godotenv.Load()

	db := config.InitDB()

	studentRepo := repository.NewStudentRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	enrollmentRepo := repository.NewEnrollmentRepository(db)

	studentHandler := handler.NewStudentHandler(studentRepo)
	courseHandler := handler.NewCourseHandler(courseRepo)
	enrollmentHandler := handler.NewEnrollmentHandler(enrollmentRepo, courseRepo)

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// Swagger route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Public routes
	e.POST("/students/register", studentHandler.Register)
	e.POST("/students/login", studentHandler.Login)

	// Protected routes (JWT)
	protected := e.Group("")
	protected.Use(middleware.JWTAuth)
	protected.GET("/students/me", studentHandler.GetMe)
	protected.GET("/courses", courseHandler.GetAll)
	protected.POST("/enrollments", enrollmentHandler.Create)
	protected.DELETE("/enrollments/:id", enrollmentHandler.Delete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
