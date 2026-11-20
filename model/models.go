package model

import "time"

type Department struct {
	DepartmentID uint   `gorm:"column:department_id;primaryKey;autoIncrement" json:"department_id"`
	Name         string `gorm:"column:name;not null" json:"name"`
	Description  string `gorm:"column:description;not null" json:"description"`
}

type Student struct {
	StudentID   uint         `gorm:"column:student_id;primaryKey;autoIncrement" json:"student_id"`
	FirstName   string       `gorm:"column:first_name;not null" json:"first_name"`
	LastName    string       `gorm:"column:last_name;not null" json:"last_name"`
	Email       string       `gorm:"column:email;unique;not null" json:"email"`
	Password    string       `gorm:"column:password;not null" json:"-"`
	Address     string       `gorm:"column:address;not null" json:"address"`
	DateOfBirth time.Time    `gorm:"column:date_of_birth;type:date;not null" json:"date_of_birth"`
	Enrollments []Enrollment `gorm:"foreignKey:StudentID;references:StudentID" json:"enrollments,omitempty"`
}

type Course struct {
	CourseID     uint       `gorm:"column:course_id;primaryKey;autoIncrement" json:"course_id"`
	Name         string     `gorm:"column:name;not null" json:"name"`
	Description  string     `gorm:"column:description;not null" json:"description"`
	DepartmentID uint       `gorm:"column:department_id;not null" json:"department_id"`
	Credits      int        `gorm:"column:credits;not null" json:"credits"`
	Department   Department `gorm:"foreignKey:DepartmentID;references:DepartmentID" json:"department,omitempty"`
}

type Enrollment struct {
	EnrollmentID   uint      `gorm:"column:enrollment_id;primaryKey;autoIncrement" json:"enrollment_id"`
	StudentID      uint      `gorm:"column:student_id;not null" json:"student_id"`
	CourseID       uint      `gorm:"column:course_id;not null" json:"course_id"`
	EnrollmentDate time.Time `gorm:"column:enrollment_date;type:date;not null" json:"enrollment_date"`
	Course         Course    `gorm:"foreignKey:CourseID;references:CourseID" json:"course,omitempty"`
}
