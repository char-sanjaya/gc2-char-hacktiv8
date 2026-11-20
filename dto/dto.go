package dto

type ErrorResponse struct {
	Message string `json:"message"`
}

type RegisterRequest struct {
	FirstName   string `json:"first_name" example:"Budi"`
	LastName    string `json:"last_name" example:"Santoso"`
	Address     string `json:"address" example:"Jl. Sudirman No. 1, Jakarta"`
	Email       string `json:"email" example:"budi@student.ac.id"`
	Password    string `json:"password" example:"rahasia123"`
	DateOfBirth string `json:"date_of_birth" example:"2001-05-14"`
}

type StudentData struct {
	StudentID   uint   `json:"student_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	DateOfBirth string `json:"date_of_birth"`
}

type RegisterResponse struct {
	Message string      `json:"message"`
	Data    StudentData `json:"data"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"budi@student.ac.id"`
	Password string `json:"password" example:"rahasia123"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type StudentEnrollmentItem struct {
	EnrollmentID   uint   `json:"enrollment_id"`
	CourseName     string `json:"course_name"`
	EnrollmentDate string `json:"enrollment_date"`
}

type StudentMeResponse struct {
	StudentID   uint                    `json:"student_id"`
	FullName    string                  `json:"full_name"`
	Email       string                  `json:"email"`
	Address     string                  `json:"address"`
	DateOfBirth string                  `json:"date_of_birth"`
	Enrollments []StudentEnrollmentItem `json:"enrollments"`
}

type CourseResponse struct {
	CourseID     uint   `json:"course_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DepartmentID uint   `json:"department_id"`
	Credits      int    `json:"credits"`
}

type EnrollmentRequest struct {
	CourseID uint `json:"course_id" example:"1"`
}

type EnrollmentData struct {
	EnrollmentID   uint           `json:"enrollment_id"`
	EnrollmentDate string         `json:"enrollment_date"`
	Course         CourseResponse `json:"course"`
}

type EnrollmentResponse struct {
	Message string         `json:"message"`
	Data    EnrollmentData `json:"data"`
}
