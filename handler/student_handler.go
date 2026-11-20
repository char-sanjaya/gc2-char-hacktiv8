package handler

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"gc2-university-api/dto"
	"gc2-university-api/middleware"
	"gc2-university-api/model"
	"gc2-university-api/repository"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type StudentHandler struct {
	Repo *repository.StudentRepository
}

func NewStudentHandler(repo *repository.StudentRepository) *StudentHandler {
	return &StudentHandler{Repo: repo}
}

// Register godoc
// @Summary Register a new student
// @Description Create a new student account
// @Tags Students
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Student registration data"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /students/register [post]
func (h *StudentHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "invalid request body"})
	}

	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Address = strings.TrimSpace(req.Address)

	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" || req.Address == "" || req.DateOfBirth == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "all fields (first_name, last_name, address, email, password, date_of_birth) are required",
		})
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "date_of_birth must be in YYYY-MM-DD format",
		})
	}

	if _, err := h.Repo.FindByEmail(req.Email); err == nil {
		return c.JSON(http.StatusConflict, dto.ErrorResponse{
			Message: "email is already registered",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "failed to process password"})
	}

	student := model.Student{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    string(hashedPassword),
		Address:     req.Address,
		DateOfBirth: dob,
	}

	if err := h.Repo.Create(&student); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "failed to register student"})
	}

	return c.JSON(http.StatusCreated, dto.RegisterResponse{
		Message: "student registered successfully",
		Data: dto.StudentData{
			StudentID:   student.StudentID,
			FirstName:   student.FirstName,
			LastName:    student.LastName,
			Email:       student.Email,
			Address:     student.Address,
			DateOfBirth: student.DateOfBirth.Format("2006-01-02"),
		},
	})
}

// Login godoc
// @Summary Student login
// @Description Authenticate student and return JWT token
// @Tags Students
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Student login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /students/login [post]
func (h *StudentHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "invalid request body"})
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "email and password are required"})
	}

	student, err := h.Repo.FindByEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "invalid email or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "invalid email or password"})
	}

	token, err := middleware.GenerateToken(student.StudentID, student.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "failed to generate token"})
	}

	return c.JSON(http.StatusOK, dto.LoginResponse{
		Message: "login successful",
		Token:   token,
	})
}

// GetMe godoc
// @Summary Get logged-in student profile
// @Description Retrieve profile details and enrolled courses of the authenticated student
// @Tags Students
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.StudentMeResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /students/me [get]
func (h *StudentHandler) GetMe(c echo.Context) error {
	studentID := c.Get("student_id").(uint)

	student, err := h.Repo.FindByIDWithEnrollments(studentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "student not found"})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "failed to retrieve student profile"})
	}

	enrollments := make([]dto.StudentEnrollmentItem, 0, len(student.Enrollments))
	for _, e := range student.Enrollments {
		enrollments = append(enrollments, dto.StudentEnrollmentItem{
			EnrollmentID:   e.EnrollmentID,
			CourseName:     e.Course.Name,
			EnrollmentDate: e.EnrollmentDate.Format("2006-01-02"),
		})
	}

	return c.JSON(http.StatusOK, dto.StudentMeResponse{
		StudentID:   student.StudentID,
		FullName:    student.FirstName + " " + student.LastName,
		Email:       student.Email,
		Address:     student.Address,
		DateOfBirth: student.DateOfBirth.Format("2006-01-02"),
		Enrollments: enrollments,
	})
}
