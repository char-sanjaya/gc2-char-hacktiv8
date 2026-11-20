package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"gc2-university-api/dto"
	"gc2-university-api/model"
	"gc2-university-api/repository"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type EnrollmentHandler struct {
	EnrollmentRepo *repository.EnrollmentRepository
	CourseRepo     *repository.CourseRepository
}

func NewEnrollmentHandler(eRepo *repository.EnrollmentRepository, cRepo *repository.CourseRepository) *EnrollmentHandler {
	return &EnrollmentHandler{
		EnrollmentRepo: eRepo,
		CourseRepo:     cRepo,
	}
}

// Create godoc
// @Summary Enroll in a course
// @Description Enroll the authenticated student into a course
// @Tags Enrollments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.EnrollmentRequest true "Course ID to enroll"
// @Success 201 {object} dto.EnrollmentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /enrollments [post]
func (h *EnrollmentHandler) Create(c echo.Context) error {
	studentID := c.Get("student_id").(uint)

	var req dto.EnrollmentRequest
	if err := c.Bind(&req); err != nil || req.CourseID == 0 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "valid course_id is required"})
	}

	course, err := h.CourseRepo.FindByID(req.CourseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "course not found"})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "failed to verify course"})
	}

	exists, err := h.EnrollmentRepo.Exists(studentID, req.CourseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "failed to check enrollment"})
	}
	if exists {
		return c.JSON(http.StatusConflict, dto.ErrorResponse{Message: "you are already enrolled in this course"})
	}

	enrollment := model.Enrollment{
		StudentID:      studentID,
		CourseID:       req.CourseID,
		EnrollmentDate: time.Now(),
	}

	if err := h.EnrollmentRepo.Create(&enrollment); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "failed to create enrollment"})
	}

	return c.JSON(http.StatusCreated, dto.EnrollmentResponse{
		Message: "enrolled successfully",
		Data: dto.EnrollmentData{
			EnrollmentID:   enrollment.EnrollmentID,
			EnrollmentDate: enrollment.EnrollmentDate.Format("2006-01-02"),
			Course: dto.CourseResponse{
				CourseID:     course.CourseID,
				Name:         course.Name,
				Description:  course.Description,
				DepartmentID: course.DepartmentID,
				Credits:      course.Credits,
			},
		},
	})
}

// Delete godoc
// @Summary Delete an enrollment
// @Description Remove a course enrollment belonging to the authenticated student
// @Tags Enrollments
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Enrollment ID"
// @Success 200 {object} dto.EnrollmentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /enrollments/{id} [delete]
func (h *EnrollmentHandler) Delete(c echo.Context) error {
	studentID := c.Get("student_id").(uint)

	idParam := c.Param("id")
	enrollmentID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || enrollmentID == 0 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "invalid enrollment id"})
	}

	enrollment, err := h.EnrollmentRepo.FindByIDAndStudent(uint(enrollmentID), studentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "enrollment not found"})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "failed to retrieve enrollment"})
	}

	if err := h.EnrollmentRepo.Delete(enrollment); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "failed to delete enrollment"})
	}

	return c.JSON(http.StatusOK, dto.EnrollmentResponse{
		Message: "enrollment deleted successfully",
		Data: dto.EnrollmentData{
			EnrollmentID:   enrollment.EnrollmentID,
			EnrollmentDate: enrollment.EnrollmentDate.Format("2006-01-02"),
			Course: dto.CourseResponse{
				CourseID:     enrollment.Course.CourseID,
				Name:         enrollment.Course.Name,
				Description:  enrollment.Course.Description,
				DepartmentID: enrollment.Course.DepartmentID,
				Credits:      enrollment.Course.Credits,
			},
		},
	})
}
