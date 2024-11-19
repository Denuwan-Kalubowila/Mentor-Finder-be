package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Denuwan-Kalubowila/mentor-finder/model"
	"github.com/Denuwan-Kalubowila/mentor-finder/utils"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	DB *sql.DB
}

// get all students
func (db *Handlers) GetStudents(c echo.Context) error {
	row, err := db.DB.Query("SELECT id,name,email,phone,university,degree FROM students")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error while fetching students"})
	}
	defer row.Close()
	var students []model.Student
	for row.Next() {
		var student model.Student
		if err := row.Scan(&student.ID, &student.Name, &student.Email, &student.Phone, &student.University, &student.Degree); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to scan students"})
		}
		students = append(students, student)
	}
	return c.JSON(http.StatusOK, students)

}

// register a student
func (db *Handlers) RegisterStudent(c echo.Context) error {
	var student model.Student
	if err := c.Bind(&student); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}
	query := `INSERT INTO students (name, email, phone, password, university, degree) 
              VALUES (?, ?, ?, ?, ?, ?)`

	hasedPass, err := utils.HashPassword(student.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
	}
	result, err := db.DB.Exec(query, student.Name, student.Email, student.Phone,
		hasedPass, student.University, student.Degree)
	if err != nil {
		log.Printf("Error registering student: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error while adding student",
		})
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error retrieving new student ID",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "student registered successfully",
		"id":      id,
	})
}
