-- 1. Tabel Departments
CREATE TABLE departments (
    department_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL
);

-- 2. Tabel Students
CREATE TABLE students (
    student_id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    date_of_birth DATE NOT NULL
);

-- 3. Tabel Professors
CREATE TABLE professors (
    professor_id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    address TEXT NOT NULL
);

-- 4. Tabel Courses
CREATE TABLE courses (
    course_id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    description TEXT NOT NULL,
    department_id INT NOT NULL REFERENCES departments(department_id) ON DELETE CASCADE,
    credits INT NOT NULL
);

-- 5. Tabel Enrollments
CREATE TABLE enrollments (
    enrollment_id SERIAL PRIMARY KEY,
    student_id INT NOT NULL REFERENCES students(student_id) ON DELETE CASCADE,
    course_id INT NOT NULL REFERENCES courses(course_id) ON DELETE CASCADE,
    enrollment_date DATE NOT NULL DEFAULT CURRENT_DATE,
    UNIQUE (student_id, course_id)
);

-- 6. Tabel Teachings
CREATE TABLE teachings (
    professor_id INT NOT NULL REFERENCES professors(professor_id) ON DELETE CASCADE,
    course_id INT NOT NULL REFERENCES courses(course_id) ON DELETE CASCADE,
    PRIMARY KEY (professor_id, course_id)
);

-- ==========================================
-- SEEDING DATA
-- ==========================================

INSERT INTO departments (name, description) VALUES
('Computer Science', 'Department focusing on software engineering, algorithms, and systems.'),
('Information Systems', 'Department focusing on business processes and enterprise IT.'),
('Data Science', 'Department focusing on statistics, machine learning, and big data.');

INSERT INTO professors (first_name, last_name, email, address) VALUES
('Budi', 'Santoso', 'budi.santoso@univ.ac.id', 'Jl. Merdeka No. 10, Jakarta'),
('Siti', 'Aminah', 'siti.aminah@univ.ac.id', 'Jl. Sudirman No. 25, Bandung'),
('Andi', 'Wijaya', 'andi.wijaya@univ.ac.id', 'Jl. Diponegoro No. 5, Surabaya');

INSERT INTO courses (name, description, department_id, credits) VALUES
('Algorithms and Data Structures', 'Fundamental algorithms and data structures in Go.', 1, 4),
('Database Systems', 'Relational database design and SQL optimization.', 1, 3),
('Enterprise Architecture', 'Designing scalable IT infrastructure for enterprises.', 2, 3),
('Machine Learning Basics', 'Introduction to supervised and unsupervised learning.', 3, 4);

INSERT INTO teachings (professor_id, course_id) VALUES
(1, 1),
(1, 2),
(2, 3),
(3, 4),
(2, 2);

-- Password dummy di bawah adalah hash bcrypt untuk 'password123'
INSERT INTO students (first_name, last_name, email, password, address, date_of_birth) VALUES
('Rina', 'Kusuma', 'rina.kusuma@student.ac.id', '$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6B6Y4p7X6eN8aE3b3f1k8d9e7g6h2', 'Jl. Melati No. 12, Jakarta', '2001-05-14'),
('Dimas', 'Pratama', 'dimas.pratama@student.ac.id', '$2a$10$7EqJtq98hPqEX7fNZaFWoOHi6B6Y4p7X6eN8aE3b3f1k8d9e7g6h2', 'Jl. Kenanga No. 8, Bekasi', '2000-11-22');

INSERT INTO enrollments (student_id, course_id, enrollment_date) VALUES
(1, 1, '2026-09-01'),
(1, 2, '2026-09-02'),
(2, 1, '2026-09-01'),
(2, 3, '2026-09-03');