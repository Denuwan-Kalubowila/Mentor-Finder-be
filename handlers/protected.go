package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ProtectedHandler(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "You are authorized",
		"user_id": userID,
	})
}
