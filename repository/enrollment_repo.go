package repository

import (
	"gc2-university-api/model"

	"gorm.io/gorm"
)

type EnrollmentRepository struct {
	DB *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{DB: db}
}

func (r *EnrollmentRepository) Exists(studentID, courseID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Enrollment{}).
		Where("student_id = ? AND course_id = ?", studentID, courseID).
		Count(&count).Error
	return count > 0, err
}

func (r *EnrollmentRepository) Create(enrollment *model.Enrollment) error {
	return r.DB.Create(enrollment).Error
}

func (r *EnrollmentRepository) FindByIDAndStudent(enrollmentID, studentID uint) (*model.Enrollment, error) {
	var enrollment model.Enrollment
	err := r.DB.Preload("Course").
		Where("enrollment_id = ? AND student_id = ?", enrollmentID, studentID).
		First(&enrollment).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *EnrollmentRepository) Delete(enrollment *model.Enrollment) error {
	return r.DB.Delete(enrollment).Error
}
