package repository

import (
	"gc2-university-api/model"

	"gorm.io/gorm"
)

type CourseRepository struct {
	DB *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{DB: db}
}

func (r *CourseRepository) FindAll() ([]model.Course, error) {
	var courses []model.Course
	err := r.DB.Order("course_id asc").Find(&courses).Error
	return courses, err
}

func (r *CourseRepository) FindByID(courseID uint) (*model.Course, error) {
	var course model.Course
	if err := r.DB.First(&course, "course_id = ?", courseID).Error; err != nil {
		return nil, err
	}
	return &course, nil
}
