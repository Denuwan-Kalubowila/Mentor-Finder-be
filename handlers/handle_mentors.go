package handlers

import (
	"log"
	"net/http"

	"github.com/Denuwan-Kalubowila/mentor-finder/model"
	"github.com/Denuwan-Kalubowila/mentor-finder/utils"
	"github.com/labstack/echo/v4"
)

// get all mentors
func (db *Handlers) GetMentors(c echo.Context) error {
	row, err := db.DB.Query("SELECT id,name,email,phone,role FROM mentors")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error while fetching mentors",
		})
	}
	defer row.Close()
	var mentors []model.Mentor
	for row.Next() {
		var mentor model.Mentor
		if err := row.Scan(&mentor.ID, &mentor.Name, &mentor.Email, &mentor.Phone, &mentor.Role); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Failed to scan mentors",
			})
		}
		mentors = append(mentors, mentor)
	}

	return c.JSON(http.StatusOK, mentors)

}

// get mentor by role
func (db *Handlers) GetMentorByRole(c echo.Context) error {
	var role struct {
		Role       string `json:"role"`
		Experience int64  `json:"experience"`
	}
	err := c.Bind(&role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}
	row, err := db.DB.Query("SELECT id,name,email,phone,role,experience FROM mentors WHERE role=? and experience=?", role.Role, role.Experience)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error while fetching mentors",
		})
	}

	defer row.Close()

	var mentors []model.Mentor
	for row.Next() {
		var mentor model.Mentor
		if err := row.Scan(&mentor.ID, &mentor.Name, &mentor.Email, &mentor.Phone, &mentor.Role, &mentor.Experience); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Failed to scan mentor",
			})
		}
		mentors = append(mentors, mentor)
	}

	return c.JSON(http.StatusOK, mentors)
}

// register a mentor
func (db *Handlers) RegisterMentor(c echo.Context) error {
	var mentor model.Mentor
	if err := c.Bind(&mentor); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}
	hasedPass, err := utils.HashPassword(mentor.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
	}
	query := `INSERT INTO mentors (name, email, phone, password, role, experience) 
              VALUES (?, ?, ?, ?, ?, ?)`

	result, err := db.DB.Exec(query, mentor.Name, mentor.Email, mentor.Phone,
		hasedPass, mentor.Role, mentor.Experience)
	if err != nil {
		log.Printf("Error registering mentor: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error while adding mentor",
		})
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error retrieving new mentor ID",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Mentor registered successfully",
		"id":      id,
	})
}

// update mentor profile
func (db *Handlers) UpdateMentor(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized nill",
		})
	}

	userId, ok := userID.(int64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user ID",
		})
	}

	var updateData struct {
		Name       string `json:"name"`
		Phone      string `json:"phone"`
		Role       string `json:"role"`
		Experience int64  `json:"experience"`
	}

	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	query := `UPDATE mentors SET name=?,phone=?,role=?,experience=? WHERE id=?`
	_, err := db.DB.Exec(query, updateData.Name, updateData.Phone, updateData.Role, updateData.Experience, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update profile"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}
