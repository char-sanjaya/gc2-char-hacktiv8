package handler

import (
	"net/http"

	"gc2-university-api/dto"
	"gc2-university-api/repository"

	"github.com/labstack/echo/v4"
)

type CourseHandler struct {
	Repo *repository.CourseRepository
}

func NewCourseHandler(repo *repository.CourseRepository) *CourseHandler {
	return &CourseHandler{Repo: repo}
}

// GetAll godoc
// @Summary Retrieve all courses
// @Description Get list of all available courses
// @Tags Courses
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} dto.CourseResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /courses [get]
func (h *CourseHandler) GetAll(c echo.Context) error {
	courses, err := h.Repo.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "failed to fetch courses"})
	}

	res := make([]dto.CourseResponse, 0, len(courses))
	for _, course := range courses {
		res = append(res, dto.CourseResponse{
			CourseID:     course.CourseID,
			Name:         course.Name,
			Description:  course.Description,
			DepartmentID: course.DepartmentID,
			Credits:      course.Credits,
		})
	}

	return c.JSON(http.StatusOK, res)
}
