package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Missing authorization token",
			})
		}

		tokenChunks := strings.Split(authHeader, " ")
		if len(tokenChunks) != 2 || tokenChunks[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Invalid authorization header format",
			})
		}

		userID, err := ValidateToken(tokenChunks[1], AccessTokenSecret)
		if err != nil {
			c.Logger().Error("Token validation failed: ", err)
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Invalid or expired token",
			})
		}

		c.Set("user_id", userID)
		return next(c)
	}
}
