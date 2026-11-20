package repository

import (
	"gc2-university-api/model"

	"gorm.io/gorm"
)

type StudentRepository struct {
	DB *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{DB: db}
}

func (r *StudentRepository) Create(student *model.Student) error {
	return r.DB.Create(student).Error
}

func (r *StudentRepository) FindByEmail(email string) (*model.Student, error) {
	var student model.Student
	if err := r.DB.Where("email = ?", email).First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *StudentRepository) FindByIDWithEnrollments(studentID uint) (*model.Student, error) {
	var student model.Student
	err := r.DB.Preload("Enrollments.Course").First(&student, "student_id = ?", studentID).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}
