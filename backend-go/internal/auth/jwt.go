package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func InitJWT(secret string) {
	jwtSecret = []byte(secret)
}

func GenerateTeacherToken(teacherID int) (string, error) {
	if len(jwtSecret) == 0 {
		return "", fmt.Errorf("jwt secret is not initialized")
	}

	claims := jwt.MapClaims{
		"sub": teacherID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseTeacherToken(tokenString string) (int, error) {
	if len(jwtSecret) == 0 {
		return 0, fmt.Errorf("jwt secret is not initialized")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	sub, ok := claims["sub"]
	if !ok {
		return 0, fmt.Errorf("sub not found")
	}

	switch v := sub.(type) {
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("invalid sub type")
	}
}