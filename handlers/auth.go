package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Denuwan-Kalubowila/mentor-finder/auth"
	"github.com/Denuwan-Kalubowila/mentor-finder/model"
	"github.com/Denuwan-Kalubowila/mentor-finder/utils"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Db *sql.DB
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (db AuthHandler) Login(c echo.Context) error {
	var request LoginRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}
	var mentor model.Mentor
	loginQuery := `SELECT id, email, password FROM mentors WHERE email = ?`
	err := db.Db.QueryRow(loginQuery, request.Email).Scan(&mentor.ID, &mentor.Email, &mentor.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Invalid email or password",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error while fetching mentor",
		})
	}

	if !utils.ComparePassword(mentor.Password, request.Password) {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Invalid email or password",
		})
	}

	accessToken, err := auth.GenarateAccessToken(mentor.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error while generating access token",
		})
	}
	refreshToken, err := auth.GenarateRefreshToken(mentor.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error while generating refresh token",
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	refreshToken := c.FormValue("refresh_token")
	if refreshToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Refresh token is required"})
	}

	userID, err := auth.ValidateToken(refreshToken, auth.RefreshTokenSecret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid refresh token"})
	}

	newAccessToken, err := auth.GenarateAccessToken(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate new access token"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token": newAccessToken,
	})
}
