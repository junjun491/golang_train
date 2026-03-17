package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("dev-secret-change-me")

func GenerateTeacherToken(teacherID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": teacherID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}