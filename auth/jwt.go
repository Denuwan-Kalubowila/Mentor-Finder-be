package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	AccessTokenSecret  = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	RefreshTokenSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
)

func GenarateAccessToken(userId int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Minute * 30).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(AccessTokenSecret)
}

func GenarateRefreshToken(userId int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(AccessTokenSecret)
}

func ValidateToken(tokenstr string, secret []byte) (int64, error) {
	token, err := jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"]
		fmt.Printf("user_id type: %T, value: %v\n", userID, userID)

		userIDNum, ok := userID.(json.Number)
		if !ok {
			return 0, fmt.Errorf("invalid user_id claim type: expected json.Number, got %T", userID)
		}

		userId, err := userIDNum.Int64()
		if err != nil {
			return 0, fmt.Errorf("failed to convert user_id to int64: %v", err)
		}

		return userId, nil
	}
	return 0, jwt.ErrSignatureInvalid
}
